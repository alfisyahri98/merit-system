package middleware

import (
	"Backend/internal/helper"
	"Backend/internal/model"
	"Backend/internal/repository"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const principalKey = "principal"

// AuthMiddleware: AUTHENTICATION — memastikan "siapa" yang request.
// 1. Ambil token dari header Authorization: Bearer <token>
// 2. Verifikasi tanda tangan & masa berlaku JWT
// 3. Lihat jenis token (kind): user → tabel users, client → tabel api_clients;
// keduanya dicek masih ada & masih aktif.
// 4. Simpan identitas (Principal) ke context request
func AuthMiddleware(secret string, userRepo repository.UserRepository, clientRepo repository.ApiClientRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			helper.Fail(c, http.StatusUnauthorized, "token tidak ditemukan", nil)
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		claims, err := helper.ParseJWT(tokenString, secret)
		if err != nil {
			helper.Fail(c, http.StatusUnauthorized, "token tidak valid atau kedaluwarsa", nil)
			return
		}

		var p *model.Principal
		switch claims.Kind {
		case model.KindUser:
			p, err = userPrincipal(userRepo, claims.UserID)
		case model.KindClient:
			p, err = clientPrincipal(clientRepo, claims.ClientID)
		default:
			err = errUnauthenticated("jenis token tidak dikenal")
		}

		var authErr authError
		if errors.As(err, &authErr) {
			helper.Fail(c, http.StatusUnauthorized, string(authErr), nil)
			return
		}
		if err != nil {
			helper.Fail(c, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		c.Set(principalKey, p)
		c.Next()
	}
}

// authError: error yang dijawab 401 (bukan 500).
type authError string

func (e authError) Error() string { return string(e) }

func errUnauthenticated(msg string) error { return authError(msg) }

// userPrincipal: role & satker diambil dari DB (bukan dari token),
// jadi perubahan role/satker/nonaktif langsung berlaku.
func userPrincipal(repo repository.UserRepository, id int) (*model.Principal, error) {
	user, err := repo.FindById(uint(id))
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, errUnauthenticated("user tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, errUnauthenticated("akun dinonaktifkan")
	}
	return &model.Principal{
		Kind:        model.KindUser,
		ID:          user.ID,
		Name:        user.Username,
		Role:        user.Role,
		SatkerID:    user.SatkerID,
		Permissions: model.RolePermissions[user.Role],
	}, nil
}

// clientPrincipal: permission = scopes, wilayah = akses_satker_id.
// Akses yang dicabut (is_active=false) langsung berlaku walau token belum expired.
func clientPrincipal(repo repository.ApiClientRepository, id int) (*model.Principal, error) {
	client, err := repo.FindByID(id)
	if errors.Is(err, repository.ErrClientNotFound) {
		return nil, errUnauthenticated("aplikasi tidak terdaftar")
	}
	if err != nil {
		return nil, err
	}
	if !client.IsActive {
		return nil, errUnauthenticated("akses aplikasi telah dicabut")
	}
	return &model.Principal{
		Kind:        model.KindClient,
		ID:          client.ID,
		Name:        client.NamaAplikasi,
		Role:        "API_CLIENT",
		SatkerID:    client.AksesSatkerID,
		Permissions: client.Scopes,
	}, nil
}

// CurrentPrincipal: ambil identitas yang disimpan AuthMiddleware.
func CurrentPrincipal(c *gin.Context) *model.Principal {
	v, ok := c.Get(principalKey)
	if !ok {
		return nil
	}
	p, _ := v.(*model.Principal)
	return p
}

// RequirePermission: AUTHORIZATION — memastikan pemanggil "boleh" melakukan aksi ini.
func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := CurrentPrincipal(c)
		if p == nil {
			helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
			return
		}
		if !p.Can(perm) {
			helper.Fail(c, http.StatusForbidden, "anda tidak memiliki akses: "+perm, nil)
			return
		}
		c.Next()
	}
}
