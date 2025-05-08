package utils

import (
	"errors"
	"time"
	"user-management/internal/config"

	"github.com/golang-jwt/jwt"
)

var accessSecretKey = []byte(config.JWTAccessSecret)
var refreshSecretKey = []byte(config.JWTRefreshSecret)

func GenerateAccessToken(userID int32, username string, isAdmin bool) (string, error) {

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecretKey)

}

func GenerateRefreshToken(userID int32, username string, isAdmin bool) (string, error) {

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecretKey)

}

func ValidateAccessToken(tokenString string) (int64, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return accessSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, false, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, errors.New("invalid token claims")
	}

	userID := int64(claims["user_id"].(float64))
	isAdmin := false
	if v, ok := claims["is_admin"].(bool); ok {
		isAdmin = v
	} else if v, ok := claims["is_admin"].(float64); ok {
		isAdmin = v == 1
	}
	return userID, isAdmin, nil
}
