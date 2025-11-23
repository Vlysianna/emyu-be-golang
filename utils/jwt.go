package utils

import (
	"errors"
	"time"

	"github.com/emyu/ecommer-be/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	RoleID      int      `json:"role_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateToken(id, email string, roleID int, roleName string, permissions []string) (string, error) {
	claims := &Claims{
		ID:          id,
		Email:       email,
		RoleID:      roleID,
		RoleName:    roleName,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTKey))
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
