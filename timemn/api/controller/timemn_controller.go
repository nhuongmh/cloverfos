package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cvfs/timemn/internal/model"
)

type T2EController struct {
	EnergyService model.EnergyMngService
}

func (pctl *T2EController) EvaluateGgSheet(gc *gin.Context) {
	// langID := gc.Param("lang-id")
	err := pctl.EnergyService.EvaluateAllFromSheet(gc, true)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, gin.H{"message": "Success"})
}
