package service

import (
	"Backend/internal/dto"
	"Backend/internal/helper"
	"Backend/internal/model"
	"Backend/internal/repository"
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// ApiClientService: pengelolaan aplikasi lain oleh Admin SSDM.
type ApiClientService struct {
	repo repository.ApiClientRepository
}

func NewApiClientService(repo repository.ApiClientRepository) *ApiClientService {
	return &ApiClientService{repo: repo}
}

// ApiClientWithSecret: dikembalikan HANYA saat create & rotate.
// Secret asli tidak disimpan di DB (cuma hash-nya), jadi tidak bisa dilihat lagi.
type ApiClientWithSecret struct {
	Client       *model.ApiClient `json:"client"`
	ClientSecret string           `json:"client_secret"`
	Peringatan   string           `json:"peringatan"`
}

const peringatanSecret = "Simpan client_secret ini sekarang. Secret tidak akan ditampilkan lagi."

func (s *ApiClientService) List(q dto.PageQuery) (*Page[model.ApiClient], error) {
	q.Normalize()
	items, total, err := s.repo.List(q.Limit, q.Offset())
	if err != nil {
		return nil, err
	}
	return &Page[model.ApiClient]{Items: items, Page: q.Page, Limit: q.Limit, Total: total}, nil
}

func (s *ApiClientService) Create(req dto.CreateApiClientRequest) (*ApiClientWithSecret, error) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}

	secret, hash, err := newSecret()
	if err != nil {
		return nil, err
	}

	client := &model.ApiClient{
		NamaAplikasi:  strings.TrimSpace(req.NamaAplikasi),
		ClientID:      req.ClientID,
		SecretHash:    hash,
		Scopes:        req.Scopes,
		AksesSatkerID: req.AksesSatkerID,
		IsActive:      true,
	}
	if err := s.repo.Create(client); err != nil {
		return nil, err
	}
	return &ApiClientWithSecret{Client: client, ClientSecret: secret, Peringatan: peringatanSecret}, nil
}

func (s *ApiClientService) Update(id int, req dto.UpdateApiClientRequest) (*model.ApiClient, error) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	client := &model.ApiClient{
		ID:            id,
		NamaAplikasi:  strings.TrimSpace(req.NamaAplikasi),
		Scopes:        req.Scopes,
		AksesSatkerID: req.AksesSatkerID,
	}
	if err := s.repo.Update(client); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

// SetStatus: false = cabut akses. Token yang sudah beredar langsung ditolak
// karena AuthMiddleware mengecek is_active setiap request.
func (s *ApiClientService) SetStatus(id int, active bool) (*model.ApiClient, error) {
	if err := s.repo.SetActive(id, active); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

// RotateSecret: ganti secret (mis. karena bocor). Secret lama langsung tidak berlaku
// untuk minta token baru; token yang sudah terbit tetap jalan sampai expired (maks 1 jam).
func (s *ApiClientService) RotateSecret(id int) (*ApiClientWithSecret, error) {
	secret, hash, err := newSecret()
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateSecret(id, hash); err != nil {
		return nil, err
	}
	client, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &ApiClientWithSecret{Client: client, ClientSecret: secret, Peringatan: peringatanSecret}, nil
}

// newSecret: 32 byte acak (crypto/rand) → base64url, lalu di-hash bcrypt.
func newSecret() (plain string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	plain = base64.RawURLEncoding.EncodeToString(b)
	hash, err = helper.HashPassword(plain)
	return plain, hash, err
}
