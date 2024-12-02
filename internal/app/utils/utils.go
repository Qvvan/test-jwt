package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qvvan/test-jwt/internal/config"
)

var SecretKey = GetSecretKey()

func GenerateAccessToken(currentIP string, userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["ip_address"] = currentIP
	claims["type_token"] = "access_token"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(SecretKey)) // Приведение SecretKey к []byte
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(currentIP, userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["ip"] = currentIP
	claims["type_token"] = "refresh_token"
	claims["exp"] = time.Now().Add(time.Hour * 96).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(SecretKey)) // Приведение SecretKey к []byte
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetSecretKey() string {
	cfg := config.MustLoad()
	return cfg.JwtSecretKey
}
