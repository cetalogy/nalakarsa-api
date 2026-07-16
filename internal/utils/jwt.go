package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type TokenPayload struct {
	Token     string
	ExpiresAt time.Time
}

// GenerateAccessToken generates a short-lived access token
func GenerateAccessToken(userID uuid.UUID, email, role, secret string, expirationSeconds int) (*TokenPayload, error) {
	expiresAt := time.Now().Add(time.Duration(expirationSeconds) * time.Second)
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "nalakarsa-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPayload{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateRefreshToken generates a long-lived refresh token
func GenerateRefreshToken(userID uuid.UUID, secret string, expirationSeconds int) (*TokenPayload, error) {
	expiresAt := time.Now().Add(time.Duration(expirationSeconds) * time.Second)
	claims := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "nalakarsa-backend",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPayload{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns its claims
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
