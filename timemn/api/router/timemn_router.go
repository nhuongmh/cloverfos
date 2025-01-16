package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cvfs/timemn/api/controller"
	"github.com/nhuongmh/cvfs/timemn/bootstrap"
	"github.com/nhuongmh/cvfs/timemn/pkg/services/energy"
)

func NewT2ERouter(app *bootstrap.T2EApplication, timeout time.Duration, publicRouter, privateRouter *gin.RouterGroup) {

	ts := energy.NewEnergyMngService(app.Env)
	tc := &controller.T2EController{EnergyService: ts}

	publicRouter.PUT(DEFAULT_API_PREFIX+"/t2e/ggsheet/evaluate", tc.EvaluateGgSheet)
	// publicRouter.POST(DEFAULT_API_PREFIX+"/practice/:card-id", tc.GetCard)

}
