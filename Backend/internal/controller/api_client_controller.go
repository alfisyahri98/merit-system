package controller

import (
	"Backend/internal/dto"
	"Backend/internal/helper"
	"Backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiClientController: endpoint Admin SSDM untuk mengelola aplikasi lain.
// Semua route dijaga RequirePermission(user:manage).
type ApiClientController struct {
	svc *service.ApiClientService
}

func NewApiClientController(svc *service.ApiClientService) *ApiClientController {
	return &ApiClientController{svc: svc}
}

// GET /api/api-clients
func (ctl *ApiClientController) List(c *gin.Context) {
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.Fail(c, http.StatusBadRequest, "parameter page/limit tidak valid", nil)
		return
	}
	items, err := ctl.svc.List(q)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, items)
}

// POST /api/api-clients
func (ctl *ApiClientController) Create(c *gin.Context) {
	var req dto.CreateApiClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}
	result, err := ctl.svc.Create(req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusCreated, result)
}

// PUT /api/api-clients/:id
func (ctl *ApiClientController) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateApiClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}
	client, err := ctl.svc.Update(int(id), req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, client)
}

// PATCH /api/api-clients/:id/status
func (ctl *ApiClientController) SetStatus(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req dto.SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "is_active wajib diisi (true/false)", nil)
		return
	}
	client, err := ctl.svc.SetStatus(int(id), *req.IsActive)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, client)
}

// POST /api/api-clients/:id/rotate-secret
func (ctl *ApiClientController) RotateSecret(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	result, err := ctl.svc.RotateSecret(int(id))
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, result)
}
