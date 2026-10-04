package service

import (
	"Backend/internal/config"
	"Backend/internal/helper"
	"Backend/internal/model"
	"Backend/internal/repository"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserInactive       = errors.New("akun nonaktif")
)

type LoginResult struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	User      *model.User `json:"user"`
}

type UserService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{userRepo: userRepo, cfg: cfg}
}

func (s *UserService) Login(username string, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindbyUsername(username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err // error DB → 500, jangan disamarkan jadi "password salah"
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	token, exp, err := helper.GenerateJWT(user, s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, ExpiresAt: exp, User: user}, nil
}
