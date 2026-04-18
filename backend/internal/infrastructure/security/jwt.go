// Package security fornece implementações de JWT e Hash (adaptadores de saída).
package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"endurance/config"
)

// Claims são os dados gravados dentro do token JWT.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct{}

// NewJWTService retorna uma implementação de JWTService.
func NewJWTService() *jwtService { return &jwtService{} }

// Generate cria um token JWT assinado com HMAC-SHA256.
func (s *jwtService) Generate(userID, role string) (string, error) {
	expiration := time.Duration(config.App.JWTExpirationHours) * time.Hour

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "endurance",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("assinar token: %w", err)
	}

	return signed, nil
}

// Validate analisa e valida o token, retornando os Claims.
func (s *jwtService) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return []byte(config.App.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
