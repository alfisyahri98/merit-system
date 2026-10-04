package dto

import (
	"Backend/internal/model"
	"regexp"
	"slices"
	"strings"
)

var clientIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{2,63}$`)

// CreateApiClientRequest: daftarkan aplikasi baru.
type CreateApiClientRequest struct {
	NamaAplikasi  string   `json:"nama_aplikasi" binding:"required,max=150"`
	ClientID      string   `json:"client_id" binding:"required"`
	Scopes        []string `json:"scopes" binding:"required,min=1"`
	AksesSatkerID *int     `json:"akses_satker_id" binding:"omitempty,gt=0"` // null = semua satker
}

func (r *CreateApiClientRequest) Validate() map[string]string {
	errs := validateClientCommon(r.NamaAplikasi, r.Scopes)
	if !clientIDPattern.MatchString(r.ClientID) {
		errs["client_id"] = "client_id 3-64 karakter: huruf kecil, angka, tanda minus (contoh: app-pmj)"
	}
	return errs
}

// UpdateApiClientRequest: ubah nama, scope, dan wilayah akses.
type UpdateApiClientRequest struct {
	NamaAplikasi  string   `json:"nama_aplikasi" binding:"required,max=150"`
	Scopes        []string `json:"scopes" binding:"required,min=1"`
	AksesSatkerID *int     `json:"akses_satker_id" binding:"omitempty,gt=0"`
}

func (r *UpdateApiClientRequest) Validate() map[string]string {
	return validateClientCommon(r.NamaAplikasi, r.Scopes)
}

// SetStatusRequest: cabut (false) / aktifkan kembali (true) akses aplikasi.
type SetStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

func validateClientCommon(nama string, scopes []string) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(nama) == "" {
		errs["nama_aplikasi"] = "nama aplikasi wajib diisi"
	}

	seen := map[string]bool{}
	for _, s := range scopes {
		if !slices.Contains(model.ValidScopes, s) {
			errs["scopes"] = "scope tidak dikenal: " + s + " (pilihan: " + strings.Join(model.ValidScopes, ", ") + ")"
			break
		}
		if seen[s] {
			errs["scopes"] = "scope duplikat: " + s
			break
		}
		seen[s] = true
	}
	return errs
}
