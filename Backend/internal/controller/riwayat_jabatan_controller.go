package controller

import (
	"Backend/internal/dto"
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RiwayatJabatanController struct {
	svc *service.RiwayatJabatanService
}

func NewRiwayatJabatanController(svc *service.RiwayatJabatanService) *RiwayatJabatanController {
	return &RiwayatJabatanController{svc: svc}
}

// GET /api/personel/:id/riwayat-jabatan
func (ctl *RiwayatJabatanController) List(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	personelID, ok := parseID(c, "id")
	if !ok {
		return
	}

	items, err := ctl.svc.List(p, personelID)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, items)
}

// POST /api/personel/:id/riwayat-jabatan
func (ctl *RiwayatJabatanController) Create(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	personelID, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req dto.RiwayatJabatanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}

	rj, err := ctl.svc.Create(p, personelID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusCreated, rj)
}

// PUT /api/riwayat-jabatan/:id
func (ctl *RiwayatJabatanController) Update(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req dto.RiwayatJabatanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}

	rj, err := ctl.svc.Update(p, id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, rj)
}

// DELETE /api/riwayat-jabatan/:id
func (ctl *RiwayatJabatanController) Delete(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := ctl.svc.Delete(p, id); err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, gin.H{"message": "riwayat jabatan berhasil dihapus"})
}
