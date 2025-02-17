package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cvfs/timemn/api/router"
	"github.com/nhuongmh/cvfs/timemn/bootstrap"
	"github.com/nhuongmh/cvfs/timemn/internal/logger"
)

func main() {
	logger.InitLog()
	app := bootstrap.Init()

	timeout := time.Duration(app.Env.ContextTimeout) * time.Second
	gine := gin.Default()

	logger.Log.Info().Msg("Setting up router...")
	router.Setup(&app, timeout, gine)
	gine.Use(cors.Default())

	logger.Log.Info().Msg("Starting server...")
	gine.Run(app.Env.ServerAddress)
}
