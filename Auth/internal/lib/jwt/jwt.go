package jwt

import (
	"auth/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	accessSecret  = []byte("1234567890") // лучше хранить в .env
	refreshSecret = []byte("1234567890") // или в config-файле
)

type Claims struct {
	UID   uuid.UUID `json:"uid"`
	Login string    `json:"login"`
	Role  string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(user models.User) (accessToken string, refreshToken string, err error) {
	now := time.Now()

	// Access Token: живет 15 минут
	accessClaims := Claims{
		UID:   user.Id,
		Login: user.Login,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token: живет 7 дней
	refreshClaims := Claims{
		UID:   user.Id,
		Login: user.Login,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
