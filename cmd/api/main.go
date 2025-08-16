package main

import (
	"gin/internal/bootstrap"
)

func main() {
	// Build and start the application using Fx
	app := bootstrap.BuildApp()
	app.Run()
}
