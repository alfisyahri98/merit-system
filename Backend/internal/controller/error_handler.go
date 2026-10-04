package controller

import (
	"Backend/internal/helper"
	"Backend/internal/repository"
	"Backend/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// respondError: satu tempat untuk menerjemahkan error → HTTP status.
func respondError(c *gin.Context, err error) {
	var vErr *service.ValidationError
	switch {
	case errors.As(err, &vErr):
		helper.Fail(c, http.StatusUnprocessableEntity, "validasi gagal", vErr.Fields)
	case errors.Is(err, repository.ErrPersonelNotFound):
		helper.Fail(c, http.StatusNotFound, "personel tidak ditemukan", nil)
	case errors.Is(err, repository.ErrRiwayatNotFound):
		helper.Fail(c, http.StatusNotFound, "riwayat jabatan tidak ditemukan", nil)
	case errors.Is(err, repository.ErrTmtSebelumJabatanAktif):
		helper.Fail(c, http.StatusUnprocessableEntity, "validasi gagal",
			map[string]string{"tmt_mulai": err.Error()})
	case errors.Is(err, repository.ErrUserNotFound):
		helper.Fail(c, http.StatusNotFound, "user tidak ditemukan", nil)
	case errors.Is(err, service.ErrSelfModify):
		helper.Fail(c, http.StatusUnprocessableEntity, err.Error(), nil)
	case errors.Is(err, repository.ErrClientNotFound):
		helper.Fail(c, http.StatusNotFound, "api client tidak ditemukan", nil)
	case errors.Is(err, service.ErrSatkerOutOfScope):
		helper.Fail(c, http.StatusForbidden, "satker di luar kewenangan anda", nil)
	default:
		if code, msg, ok := helper.MapPgError(err); ok {
			helper.Fail(c, code, msg, nil)
			return
		}
		helper.Fail(c, http.StatusInternalServerError, "internal server error", nil)
	}
}

// parseID: ambil :id dari URL, harus angka positif.
func parseID(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		helper.Fail(c, http.StatusBadRequest, name+" harus berupa angka positif", nil)
		return 0, false
	}
	return id, true
}
