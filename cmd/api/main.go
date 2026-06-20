package main

import "gin/internal/infra/bootstrap"

// @title Gin Skeleton API
// @version 1.0
// @description Starter Gin API with JWT auth, PostgreSQL, Fx dependency injection, and Swagger docs.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT access token. Format: "Bearer {accessToken}".
func main() {
	app := bootstrap.BuildApp()
	app.Run()
}
