package dto

import (
	"Backend/internal/model"
	"regexp"
)

// Sama dengan CHECK constraint kolom users.username
var usernamePattern = regexp.MustCompile(`^[a-z0-9_.]{4,50}$`)

const minPasswordLen = 8

// CreateUserRequest: admin membuat akun baru.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=ADMIN_SSDM OPERATOR"`
	SatkerID *int   `json:"satker_id" binding:"omitempty,gt=0"`
}

func (r *CreateUserRequest) Validate() map[string]string {
	errs := validateRoleSatker(r.Role, r.SatkerID)
	if !usernamePattern.MatchString(r.Username) {
		errs["username"] = "username 4-50 karakter: huruf kecil, angka, titik, garis bawah (contoh: opr.jaktim)"
	}
	if len(r.Password) < minPasswordLen {
		errs["password"] = "password minimal 8 karakter"
	}
	return errs
}

// UpdateUserRequest: ganti role / pindah satker.
type UpdateUserRequest struct {
	Role     string `json:"role" binding:"required,oneof=ADMIN_SSDM OPERATOR"`
	SatkerID *int   `json:"satker_id" binding:"omitempty,gt=0"`
}

func (r *UpdateUserRequest) Validate() map[string]string {
	return validateRoleSatker(r.Role, r.SatkerID)
}

// ResetPasswordRequest: admin set password baru untuk user.
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (r *ResetPasswordRequest) Validate() map[string]string {
	errs := map[string]string{}
	if len(r.Password) < minPasswordLen {
		errs["password"] = "password minimal 8 karakter"
	}
	return errs
}

// Operator WAJIB terikat satker (itu batas wilayahnya).
// Admin TIDAK BOLEH punya satker — kalau diisi, aksesnya ikut terbatas ke satker itu.
func validateRoleSatker(role string, satkerID *int) map[string]string {
	errs := map[string]string{}
	if role == model.RoleOperator && satkerID == nil {
		errs["satker_id"] = "operator wajib memiliki satker"
	}
	if role == model.RoleAdminSSDM && satkerID != nil {
		errs["satker_id"] = "admin SSDM tidak terikat satker, kosongkan satker_id"
	}
	return errs
}
