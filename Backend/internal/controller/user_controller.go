package controller

import (
	"Backend/internal/dto"
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController: menu "Kelola Pengguna" — khusus Admin SSDM (user:manage).
type UserController struct {
	svc *service.UserAdminService
}

func NewUserController(svc *service.UserAdminService) *UserController {
	return &UserController{svc: svc}
}

// GET /api/users
func (ctl *UserController) List(c *gin.Context) {
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.Fail(c, http.StatusBadRequest, "parameter page/limit tidak valid", nil)
		return
	}
	users, err := ctl.svc.List(q)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, users)
}

// POST /api/users
func (ctl *UserController) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}
	user, err := ctl.svc.Create(req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusCreated, user)
}

// PUT /api/users/:id
func (ctl *UserController) Update(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "format data tidak valid", err.Error())
		return
	}
	user, err := ctl.svc.Update(p, int(id), req)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, user)
}

// PATCH /api/users/:id/status
func (ctl *UserController) SetStatus(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req dto.SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "is_active wajib diisi (true/false)", nil)
		return
	}
	user, err := ctl.svc.SetStatus(p, int(id), *req.IsActive)
	if err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, user)
}

// POST /api/users/:id/reset-password
func (ctl *UserController) ResetPassword(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusUnprocessableEntity, "password wajib diisi", nil)
		return
	}
	if err := ctl.svc.ResetPassword(int(id), req); err != nil {
		respondError(c, err)
		return
	}
	helper.OK(c, http.StatusOK, gin.H{"message": "password berhasil direset"})
}
