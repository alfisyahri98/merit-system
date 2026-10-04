package controller

import (
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReferensiController struct {
	repo repository.ReferensiRepository
}

func NewReferensiController(repo repository.ReferensiRepository) *ReferensiController {
	return &ReferensiController{repo: repo}
}

func reply[T any](c *gin.Context, items []T, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, items)
}

// GET /api/referensi/pangkat
func (ctl *ReferensiController) Pangkat(c *gin.Context) {
	items, err := ctl.repo.Pangkat()
	reply(c, items, err)
}

// GET /api/referensi/satker
// Default: satker di wilayah pemanggil (untuk form personel).
// ?semua=1: seluruh satker (untuk riwayat jabatan, yang bisa di satker lain sebelum mutasi).
func (ctl *ReferensiController) Satker(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	scope := p.SatkerID
	if c.Query("semua") == "1" {
		scope = nil
	}
	items, err := ctl.repo.Satker(scope)
	reply(c, items, err)
}

// GET /api/referensi/fungsi
func (ctl *ReferensiController) Fungsi(c *gin.Context) {
	items, err := ctl.repo.Fungsi()
	reply(c, items, err)
}

// GET /api/referensi/nivelering
func (ctl *ReferensiController) Nivelering(c *gin.Context) {
	items, err := ctl.repo.Nivelering()
	reply(c, items, err)
}
