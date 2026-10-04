package service

import (
	"Backend/internal/config"
	"Backend/internal/helper"
	"Backend/internal/repository"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidClient  = errors.New("client_id atau client_secret salah")
	ErrClientInactive = errors.New("akses aplikasi telah dicabut")
)

// Token aplikasi sengaja dibuat pendek: kalau bocor, cepat kedaluwarsa.
const clientTokenTTL = time.Hour

type ClientTokenResult struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"` // detik
	ExpiresAt   time.Time `json:"expires_at"`
	Scopes      []string  `json:"scopes"`
}

type ClientAuthService struct {
	repo repository.ApiClientRepository
	cfg  *config.Config
}

func NewClientAuthService(repo repository.ApiClientRepository, cfg *config.Config) *ClientAuthService {
	return &ClientAuthService{repo: repo, cfg: cfg}
}

// IssueToken: OAuth2 client credentials — tukar client_id + secret dengan access token.
func (s *ClientAuthService) IssueToken(clientID, clientSecret string) (*ClientTokenResult, error) {
	client, err := s.repo.FindByClientID(clientID)
	if errors.Is(err, repository.ErrClientNotFound) {
		return nil, ErrInvalidClient
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.SecretHash), []byte(clientSecret)); err != nil {
		return nil, ErrInvalidClient
	}
	if !client.IsActive {
		return nil, ErrClientInactive
	}

	token, exp, err := helper.GenerateClientJWT(client, s.cfg.JWTSecret, clientTokenTTL)
	if err != nil {
		return nil, err
	}
	return &ClientTokenResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(clientTokenTTL.Seconds()),
		ExpiresAt:   exp,
		Scopes:      client.Scopes,
	}, nil
}
