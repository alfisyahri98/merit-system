package controller

import (
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ctl *AuthController) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		helper.Fail(c, http.StatusBadRequest, "format input tidak valid", nil)
		return
	}

	result, err := ctl.userService.Login(request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			helper.Fail(c, http.StatusUnauthorized, "username atau password salah", nil)
		case errors.Is(err, service.ErrUserInactive):
			helper.Fail(c, http.StatusForbidden, "user is inactive", nil)
		default:
			helper.Fail(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	helper.OK(c, http.StatusOK, result)
}

func (ctl *AuthController) Me(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	helper.OK(c, http.StatusOK, gin.H{
		"kind":        p.Kind,
		"id":          p.ID,
		"name":        p.Name,
		"role":        p.Role,
		"satker_id":   p.SatkerID,
		"permissions": p.Permissions,
	})
}
