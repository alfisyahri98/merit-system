package controller

import (
	"Backend/internal/dto"
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/report"
	"Backend/internal/repository"
	"Backend/internal/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PersonelController struct {
	svc *service.PersonelService
}

func NewPersonelController(svc *service.PersonelService) *PersonelController {
	return &PersonelController{svc: svc}
}

// Query string GET /personel — divalidasi otomatis oleh Gin
type listPersonelQuery struct {
	Q            string `form:"q" binding:"max=100"`
	SatkerID     *int   `form:"satker_id" binding:"omitempty,gt=0"`
	Status       string `form:"status" binding:"omitempty,oneof=AKTIF PENSIUN MENINGGAL DIBERHENTIKAN MUTASI_KELUAR"`
	Jenis        string `form:"jenis" binding:"omitempty,oneof=POLRI PNS"`
	TanpaJabatan bool   `form:"tanpa_jabatan"`
	Kelompok     string `form:"kelompok" binding:"max=50"`
	AkanPensiun  bool   `form:"akan_pensiun"`
	Sort         string `form:"sort" binding:"omitempty,oneof=nama -nama pangkat -pangkat satker -satker pensiun -pensiun"`
	Page         int    `form:"page" binding:"omitempty,min=1"`
	Limit        int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// GET /api/personel
func (ctl *PersonelController) List(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}

	var q listPersonelQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.Fail(c, http.StatusBadRequest, "parameter tidak valid", err.Error())
		return
	}

	result, err := ctl.svc.List(p, repository.PersonelFilter{
		Q:            q.Q,
		SatkerID:     q.SatkerID,
		Status:       q.Status,
		Jenis:        q.Jenis,
		TanpaJabatan: q.TanpaJabatan,
		Kelompok:     q.Kelompok,
		AkanPensiun:  q.AkanPensiun,
		Urut:         q.Sort,
		Page:         q.Page,
		Limit:        q.Limit,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, result)
}

// GET /api/personel/:id
func (ctl *PersonelController) GetByID(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	personel, err := ctl.svc.GetByID(p, id)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, personel)
}

// GET /api/personel/:id/drh → PDF Daftar Riwayat Hidup
func (ctl *PersonelController) DRH(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	personel, err := ctl.svc.GetByID(p, id) // scope satker tetap berlaku
	if err != nil {
		respondError(c, err)
		return
	}
	file, err := report.DaftarRiwayatHidup(personel, p.Name, time.Now())
	if err != nil {
		helper.Fail(c, http.StatusInternalServerError, "gagal membuat PDF", nil)
		return
	}
	nama := fmt.Sprintf("DRH_%s_%s.pdf", personel.NrpNip, strings.ReplaceAll(strings.ToUpper(personel.Nama), " ", "_"))
	c.Header("Content-Disposition", `attachment; filename="`+nama+`"`)
	c.Data(http.StatusOK, "application/pdf", file)
}

// POST /api/personel
func (ctl *PersonelController) Create(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}

	var req dto.PersonelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}

	personel, err := ctl.svc.Create(p, req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusCreated, personel)
}

// PUT /api/personel/:id
func (ctl *PersonelController) Update(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req dto.PersonelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}

	personel, err := ctl.svc.Update(p, id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, personel)
}

// DELETE /api/personel/:id
func (ctl *PersonelController) Delete(c *gin.Context) {
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
	helper.OK(c, http.StatusOK, gin.H{"message": "personel berhasil dihapus"})
}
