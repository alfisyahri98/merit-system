package controller

import (
	"Backend/internal/helper"
	"Backend/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	svc *service.ClientAuthService
}

func NewOAuthController(svc *service.ClientAuthService) *OAuthController {
	return &OAuthController{svc: svc}
}

type tokenRequest struct {
	GrantType    string `json:"grant_type" binding:"omitempty,eq=client_credentials"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
}

// POST /api/oauth/token
func (ctl *OAuthController) Token(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Fail(c, http.StatusBadRequest, "client_id dan client_secret wajib diisi", nil)
		return
	}

	result, err := ctl.svc.IssueToken(req.ClientID, req.ClientSecret)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidClient):
			helper.Fail(c, http.StatusUnauthorized, err.Error(), nil)
		case errors.Is(err, service.ErrClientInactive):
			helper.Fail(c, http.StatusForbidden, err.Error(), nil)
		default:
			helper.Fail(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}
	helper.OK(c, http.StatusOK, result)
}
