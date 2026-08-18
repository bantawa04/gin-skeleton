package router

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authhandler "gin/internal/domain/auth/handler"
	healthhandler "gin/internal/domain/health/handler"
	refreshtoken "gin/internal/domain/refresh_token"
	userdomain "gin/internal/domain/user"
	userhandler "gin/internal/domain/user/handler"
	"gin/internal/infra/logger"
	"gin/internal/shared/constant"
	exceptions "gin/internal/shared/exception"
	"gin/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testJWTSecret = "test-secret-key-for-api-tests"

type stubDriver struct{}

type stubConn struct{}

type stubTx struct{}

func (stubDriver) Open(string) (driver.Conn, error) { return &stubConn{}, nil }
func (*stubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not supported by the test database")
}
func (*stubConn) Close() error               { return nil }
func (*stubConn) Begin() (driver.Tx, error)  { return &stubTx{}, nil }
func (*stubConn) Ping(context.Context) error { return nil }
func (*stubTx) Commit() error                { return nil }
func (*stubTx) Rollback() error              { return nil }

type fakeUserService struct {
	getAllUsersPaginatedFn func(context.Context, int, int) ([]*userdomain.User, int64, error)
	getUserByIDFn          func(context.Context, string) (*userdomain.User, error)
	createUserFn           func(context.Context, userdomain.SignupInput) (*userdomain.User, error)
	updateUserFn           func(context.Context, map[string]interface{}, *string, string) (*userdomain.User, error)
	deleteUserFn           func(context.Context, string) error
	getUserByEmailFn       func(context.Context, string) (*userdomain.User, error)
}

func (f *fakeUserService) GetAllUsers(context.Context) ([]*userdomain.User, error) {
	return nil, nil
}

func (f *fakeUserService) GetAllUsersPaginated(ctx context.Context, page, perPage int) ([]*userdomain.User, int64, error) {
	if f.getAllUsersPaginatedFn != nil {
		return f.getAllUsersPaginatedFn(ctx, page, perPage)
	}
	return nil, 0, nil
}

func (f *fakeUserService) GetUserByID(ctx context.Context, id string) (*userdomain.User, error) {
	if f.getUserByIDFn != nil {
		return f.getUserByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeUserService) CreateUser(ctx context.Context, input userdomain.SignupInput) (*userdomain.User, error) {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, input)
	}
	return &userdomain.User{}, nil
}

func (f *fakeUserService) UpdateUser(ctx context.Context, updates map[string]interface{}, password *string, id string) (*userdomain.User, error) {
	if f.updateUserFn != nil {
		return f.updateUserFn(ctx, updates, password, id)
	}
	return nil, nil
}

func (f *fakeUserService) DeleteUser(ctx context.Context, id string) error {
	if f.deleteUserFn != nil {
		return f.deleteUserFn(ctx, id)
	}
	return nil
}

func (f *fakeUserService) GetUserByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	if f.getUserByEmailFn != nil {
		return f.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}

type fakeRefreshTokenService struct {
	createFn              func(context.Context, *refreshtoken.RefreshToken) (*refreshtoken.RefreshToken, error)
	findByTokenFn         func(context.Context, string) (*refreshtoken.RefreshToken, error)
	revokeByTokenFn       func(context.Context, string) error
	revokeAllUserTokensFn func(context.Context, string) error
}

func (f *fakeRefreshTokenService) Create(ctx context.Context, token *refreshtoken.RefreshToken) (*refreshtoken.RefreshToken, error) {
	if f.createFn != nil {
		return f.createFn(ctx, token)
	}
	return token, nil
}

func (f *fakeRefreshTokenService) FindByToken(ctx context.Context, token string) (*refreshtoken.RefreshToken, error) {
	if f.findByTokenFn != nil {
		return f.findByTokenFn(ctx, token)
	}
	return nil, nil
}

func (f *fakeRefreshTokenService) RevokeByToken(ctx context.Context, token string) error {
	if f.revokeByTokenFn != nil {
		return f.revokeByTokenFn(ctx, token)
	}
	return nil
}

func (f *fakeRefreshTokenService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	if f.revokeAllUserTokensFn != nil {
		return f.revokeAllUserTokensFn(ctx, userID)
	}
	return nil
}

func newTestRouter(t *testing.T, users *fakeUserService, refreshTokens *fakeRefreshTokenService) (*gin.Engine, *utils.JWTManager) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	logger.Logger = logrus.New()
	logger.ErrorLogger = logrus.New()
	logger.Logger.SetOutput(io.Discard)
	logger.ErrorLogger.SetOutput(io.Discard)

	driverName := "gin-skeleton-api-test-" + time.Now().Format("150405.000000000")
	sql.Register(driverName, stubDriver{})
	sqlDB, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("create gorm test database: %v", err)
	}

	jwtManager := utils.NewJWTManager(testJWTSecret, 15*time.Minute, 24*time.Hour)
	userHandler := userhandler.NewUserHandler(users)
	authHandler := authhandler.NewAuthHandler(users, jwtManager, refreshTokens)
	healthHandler := healthhandler.NewHealthHandler(db)

	engine := gin.New()
	engine.Use(exceptions.ErrorHandler())

	api := engine.Group("/api")
	registerWebRoutes(api, &routerDeps{
		userHandler:   userHandler,
		authHandler:   authHandler,
		healthHandler: healthHandler,
		jwtManager:    jwtManager,
		db:            db,
	})

	return engine, jwtManager
}

func performJSONRequest(t *testing.T, engine http.Handler, method, path string, body interface{}, accessToken string) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		requestBody = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	return recorder
}

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, want int) {
	t.Helper()
	if recorder.Code != want {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, want, recorder.Body.String())
	}
}

func assertSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()
	var body struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, recorder.Body.String())
	}
	if !body.Success {
		t.Fatalf("expected success response; body=%s", recorder.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	engine, _ := newTestRouter(t, &fakeUserService{}, &fakeRefreshTokenService{})

	response := performJSONRequest(t, engine, http.MethodGet, "/api/health", nil, "")

	assertStatus(t, response, http.StatusOK)
	assertSuccessResponse(t, response)
}

func TestSignupEndpoint(t *testing.T) {
	var received userdomain.SignupInput
	users := &fakeUserService{
		createUserFn: func(_ context.Context, input userdomain.SignupInput) (*userdomain.User, error) {
			received = input
			return &userdomain.User{ID: "user-1", Email: input.Email}, nil
		},
	}
	engine, _ := newTestRouter(t, users, &fakeRefreshTokenService{})

	response := performJSONRequest(t, engine, http.MethodPost, "/api/auth/signup", map[string]string{
		"first_name": "Test",
		"last_name":  "User",
		"email":      "test@example.com",
	}, "")

	assertStatus(t, response, http.StatusCreated)
	assertSuccessResponse(t, response)
	if received.Email != "test@example.com" || received.FirstName != "Test" || received.LastName != "User" {
		t.Fatalf("unexpected signup input: %+v", received)
	}
}

func TestLoginEndpoint(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	users := &fakeUserService{
		getUserByEmailFn: func(_ context.Context, email string) (*userdomain.User, error) {
			return &userdomain.User{
				ID:       "user-1",
				Email:    email,
				Password: string(passwordHash),
				Status:   constant.UserStatusActive,
			}, nil
		},
	}
	var savedRefreshToken *refreshtoken.RefreshToken
	refreshTokens := &fakeRefreshTokenService{
		createFn: func(_ context.Context, token *refreshtoken.RefreshToken) (*refreshtoken.RefreshToken, error) {
			savedRefreshToken = token
			return token, nil
		},
	}
	engine, _ := newTestRouter(t, users, refreshTokens)

	response := performJSONRequest(t, engine, http.MethodPost, "/api/auth/login", map[string]string{
		"email":    "test@example.com",
		"password": "secret123",
	}, "")

	assertStatus(t, response, http.StatusOK)
	assertSuccessResponse(t, response)
	if savedRefreshToken == nil || savedRefreshToken.UserID != "user-1" || savedRefreshToken.Token == "" {
		t.Fatalf("refresh token was not persisted correctly: %+v", savedRefreshToken)
	}
}

func TestRefreshEndpoint(t *testing.T) {
	users := &fakeUserService{}
	refreshTokens := &fakeRefreshTokenService{}
	engine, jwtManager := newTestRouter(t, users, refreshTokens)

	oldToken, err := jwtManager.GenerateRefreshToken("user-1")
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}

	revoked := false
	created := false
	refreshTokens.findByTokenFn = func(_ context.Context, token string) (*refreshtoken.RefreshToken, error) {
		if token != oldToken {
			t.Fatalf("FindByToken token = %q, want generated refresh token", token)
		}
		return &refreshtoken.RefreshToken{
			UserID:    "user-1",
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
		}, nil
	}
	refreshTokens.revokeByTokenFn = func(_ context.Context, token string) error {
		revoked = token == oldToken
		return nil
	}
	refreshTokens.createFn = func(_ context.Context, token *refreshtoken.RefreshToken) (*refreshtoken.RefreshToken, error) {
		created = token.Token != "" && token.UserID == "user-1"
		return token, nil
	}

	response := performJSONRequest(t, engine, http.MethodPost, "/api/auth/refresh", map[string]string{
		"refresh_token": oldToken,
	}, "")

	assertStatus(t, response, http.StatusOK)
	assertSuccessResponse(t, response)
	if !revoked || !created {
		t.Fatalf("refresh rotation incomplete: revoked=%v created=%v", revoked, created)
	}
}

func TestProtectedUserEndpointRejectsMissingToken(t *testing.T) {
	updateCalled := false
	users := &fakeUserService{
		updateUserFn: func(context.Context, map[string]interface{}, *string, string) (*userdomain.User, error) {
			updateCalled = true
			return &userdomain.User{}, nil
		},
	}
	engine, _ := newTestRouter(t, users, &fakeRefreshTokenService{})

	response := performJSONRequest(t, engine, http.MethodPut, "/api/users/user-1", map[string]string{
		"name": "Updated User",
	}, "")

	assertStatus(t, response, http.StatusUnauthorized)
	if updateCalled {
		t.Fatal("protected handler was called without an access token")
	}
}

func TestUserEndpoints(t *testing.T) {
	users := &fakeUserService{
		getAllUsersPaginatedFn: func(_ context.Context, page, perPage int) ([]*userdomain.User, int64, error) {
			if page != 2 || perPage != 5 {
				t.Fatalf("pagination = (%d, %d), want (2, 5)", page, perPage)
			}
			return []*userdomain.User{{ID: "user-1", Email: "test@example.com"}}, 1, nil
		},
		getUserByIDFn: func(_ context.Context, id string) (*userdomain.User, error) {
			return &userdomain.User{ID: id, Email: "test@example.com"}, nil
		},
		updateUserFn: func(_ context.Context, updates map[string]interface{}, _ *string, id string) (*userdomain.User, error) {
			if id != "user-1" || updates["first_name"] != "Updated User" {
				t.Fatalf("unexpected update: id=%q updates=%v", id, updates)
			}
			name := "Updated User"
			return &userdomain.User{ID: id, FirstName: &name, Email: "test@example.com"}, nil
		},
		deleteUserFn: func(_ context.Context, id string) error {
			if id != "user-1" {
				t.Fatalf("delete id = %q, want user-1", id)
			}
			return nil
		},
	}
	engine, jwtManager := newTestRouter(t, users, &fakeRefreshTokenService{})
	accessToken, err := jwtManager.GenerateAccessToken("user-1")
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   interface{}
		token  string
	}{
		{name: "list", method: http.MethodGet, path: "/api/users?page=2&per_page=5"},
		{name: "get by id", method: http.MethodGet, path: "/api/users/user-1"},
		{name: "update", method: http.MethodPut, path: "/api/users/user-1", body: map[string]string{"name": "Updated User"}, token: accessToken},
		{name: "delete", method: http.MethodDelete, path: "/api/users/user-1", token: accessToken},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := performJSONRequest(t, engine, test.method, test.path, test.body, test.token)
			assertStatus(t, response, http.StatusOK)
			assertSuccessResponse(t, response)
		})
	}
}

func TestLogoutEndpoint(t *testing.T) {
	var revokedUserID string
	refreshTokens := &fakeRefreshTokenService{
		revokeAllUserTokensFn: func(_ context.Context, userID string) error {
			revokedUserID = userID
			return nil
		},
	}
	engine, jwtManager := newTestRouter(t, &fakeUserService{}, refreshTokens)
	accessToken, err := jwtManager.GenerateAccessToken("user-1")
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	response := performJSONRequest(t, engine, http.MethodPost, "/api/auth/logout", nil, accessToken)

	assertStatus(t, response, http.StatusOK)
	assertSuccessResponse(t, response)
	if revokedUserID != "user-1" {
		t.Fatalf("revoked user id = %q, want user-1", revokedUserID)
	}
}
