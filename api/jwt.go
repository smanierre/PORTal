package api

import (
	"PORTal/types"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Admin bool `json:"admin"`
}

func (s Server) jwtKeyFunc(t *jwt.Token) (interface{}, error) {
	return []byte(s.config.JWTSecret), nil
}

func createToken(member types.Member, expiration time.Duration, key []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test",
			Subject:   member.ID,
			Audience:  []string{"test"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
		Admin: member.Admin,
	})
	signedToken, err := t.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func validateToken(token string, keyFunc jwt.Keyfunc, logger *slog.Logger) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, keyFunc)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelInfo, "Error validating token", slog.String("error", err.Error()))
		return nil, err
	}
	return t, err
}
