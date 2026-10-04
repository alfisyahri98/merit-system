package service

import (
	"Backend/internal/dto"
	"Backend/internal/helper"
	"Backend/internal/model"
	"Backend/internal/repository"
	"errors"
)

// Admin tidak boleh menurunkan role / menonaktifkan akunnya sendiri,
// supaya tidak ada kondisi sistem tanpa admin karena salah klik.
var ErrSelfModify = errors.New("tidak bisa mengubah role atau menonaktifkan akun sendiri")

// UserAdminService: pengelolaan akun pengguna oleh Admin SSDM.
type UserAdminService struct {
	repo repository.UserRepository
}

func NewUserAdminService(repo repository.UserRepository) *UserAdminService {
	return &UserAdminService{repo: repo}
}

func (s *UserAdminService) List(q dto.PageQuery) (*Page[model.User], error) {
	q.Normalize()
	items, total, err := s.repo.List(q.Limit, q.Offset())
	if err != nil {
		return nil, err
	}
	return &Page[model.User]{Items: items, Page: q.Page, Limit: q.Limit, Total: total}, nil
}

func (s *UserAdminService) Create(req dto.CreateUserRequest) (*model.User, error) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Role:         req.Role,
		SatkerID:     req.SatkerID,
		IsActive:     true,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Update: ganti role / pindah satker. Berlaku di request berikutnya
// karena AuthMiddleware membaca role & satker dari DB setiap request.
func (s *UserAdminService) Update(p *model.Principal, id int, req dto.UpdateUserRequest) (*model.User, error) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	if p.ID == id && req.Role != p.Role {
		return nil, ErrSelfModify
	}
	if err := s.repo.UpdateRoleSatker(id, req.Role, req.SatkerID); err != nil {
		return nil, err
	}
	return s.repo.FindById(uint(id))
}

// SetStatus: false = akun tidak bisa login, token yang beredar langsung ditolak.
func (s *UserAdminService) SetStatus(p *model.Principal, id int, active bool) (*model.User, error) {
	if p.ID == id && !active {
		return nil, ErrSelfModify
	}
	if err := s.repo.SetActive(id, active); err != nil {
		return nil, err
	}
	return s.repo.FindById(uint(id))
}

func (s *UserAdminService) ResetPassword(id int, req dto.ResetPasswordRequest) error {
	if errs := req.Validate(); len(errs) > 0 {
		return &ValidationError{Fields: errs}
	}
	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(id, hash)
}
