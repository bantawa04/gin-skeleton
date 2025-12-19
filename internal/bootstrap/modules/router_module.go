package modules

import (
	"gin/internal/config"
	"gin/internal/router"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// RouterModule provides router and HTTP server dependencies
var RouterModule = fx.Options(
	fx.Provide(router.NewRouter),
	fx.Provide(newHTTPServer),
)

// newHTTPServer creates an HTTP server with the provided router and configuration
func newHTTPServer(router *gin.Engine, cfg *config.Config) *http.Server {
	serverConfig := cfg.Server()
	return &http.Server{
		Addr:         "0.0.0.0:" + serverConfig.Port,
		Handler:      router,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
	}
}
