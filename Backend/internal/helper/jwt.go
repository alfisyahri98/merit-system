package helper

import (
	"Backend/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims: isi token. Kind membedakan token user vs token aplikasi (client).
type JWTClaims struct {
	Kind     string `json:"kind"`
	UserID   int    `json:"uid,omitempty"`
	ClientID int    `json:"cid,omitempty"`
	Role     string `json:"role,omitempty"`
	SatkerID *int   `json:"satker_id,omitempty"`
	jwt.RegisteredClaims
}

func sign(claims JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateJWT: token untuk user (login username/password).
func GenerateJWT(user *model.User, secret string, expireHours int) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(time.Duration(expireHours) * time.Hour)
	tokenString, err := sign(JWTClaims{
		Kind:     model.KindUser,
		UserID:   user.ID,
		Role:     user.Role,
		SatkerID: user.SatkerID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}, secret)
	return tokenString, exp, err
}

// GenerateClientJWT: token untuk aplikasi lain (client credentials).
func GenerateClientJWT(client *model.ApiClient, secret string, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)
	tokenString, err := sign(JWTClaims{
		Kind:     model.KindClient,
		ClientID: client.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   client.ClientID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}, secret)
	return tokenString, exp, err
}

func ParseJWT(tokenString string, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})) // tolak algoritma selain HS256
	if err != nil {
		return nil, err
	}
	if claims.Kind == "" {
		claims.Kind = model.KindUser // token lama sebelum ada field kind
	}
	return claims, nil
}
