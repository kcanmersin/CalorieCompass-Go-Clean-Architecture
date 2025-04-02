package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"CalorieCompass/internal/entity"
)

type JWTRepo struct {
	secretKey     string
	expirationHrs int
}

func NewJWTRepo(secretKey string, expirationHrs int) *JWTRepo {
	return &JWTRepo{
		secretKey:     secretKey,
		expirationHrs: expirationHrs,
	}
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (r *JWTRepo) GenerateToken(user entity.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(r.expirationHrs) * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(r.secretKey))
	if err != nil {
		return "", fmt.Errorf("generate token error: %w", err)
	}

	return tokenString, nil
}

func (r *JWTRepo) ValidateToken(tokenStr string) (int64, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(r.secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("parse token error: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}