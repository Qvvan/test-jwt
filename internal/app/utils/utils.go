package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qvvan/test-jwt/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = []byte(GetSecretKey())

const (
	TokenLength = 64
)

func GenerateAccessToken(currentIP string, userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["ip_address"] = currentIP
	claims["type_token"] = "access_token"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID string, ip string) (string, string) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, 32)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	hashString := HashToken(string(randomString))

	data := fmt.Sprintf("%s|%s|%s", userID, ip, hashString)

	encoded := base64.URLEncoding.EncodeToString([]byte(data))

	return hashString, encoded
}

func DecodeUserData(encodedString string) (string, string, string, error) {
	decoded, err := base64.URLEncoding.DecodeString(encodedString)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to decode string: %v", err)
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid string format")
	}

	return parts[0], parts[1], parts[2], nil
}

func HashToken(token string) string {
	shaHash := sha512.Sum512([]byte(token))

	bcryptHash, err := bcrypt.GenerateFromPassword(shaHash[:], bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(bcryptHash)
}

func GetSecretKey() string {
	cfg := config.MustLoad()
	return cfg.JwtSecretKey
}
