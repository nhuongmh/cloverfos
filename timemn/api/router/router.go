package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cvfs/timemn/bootstrap"
)

const (
	API_V1             = "v1"
	DEFAULT_API_PREFIX = "/api/" + API_V1
)

func Setup(app *bootstrap.T2EApplication, timeout time.Duration, gine *gin.Engine) {
	publicRouter := gine.Group("public")
	// privateRounter := gine.Group("private")

	publicRouter.Use(cors.Default())
	publicRouter.Static("/data", "./data")
	NewT2ERouter(app, timeout, publicRouter, publicRouter)

}
