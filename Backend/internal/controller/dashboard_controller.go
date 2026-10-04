package controller

import (
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	repo repository.DashboardRepository
}

func NewDashboardController(repo repository.DashboardRepository) *DashboardController {
	return &DashboardController{repo: repo}
}

// GET /api/dashboard
func (ctl *DashboardController) Ringkasan(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	data, err := ctl.repo.Ringkasan(p.SatkerID)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, data)
}
