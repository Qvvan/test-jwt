// internal/app/api/v1/manager.go
package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/qvvan/test-jwt/internal/app/utils"
	errorDb "github.com/qvvan/test-jwt/pkg/client/postgresql/utils"
	"golang.org/x/crypto/bcrypt"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (m *Manager) RefreshTokens(c *gin.Context) {
	resp, err := m.RefreshTokensService(c)
	if err != nil {
		var pubErr *PublicError
		switch {
		case errors.As(err, &pubErr):
			c.JSON(pubErr.status, gin.H{"error": pubErr.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (m *Manager) RefreshTokensService(c *gin.Context) (*RefreshResponse, error) {
	var req RefreshRequest

	// Парсим входящие данные
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, NewPublicErr(err, http.StatusBadRequest)
	}

	currentIP := c.ClientIP()

	// Парсим refresh токен
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(utils.GetSecretKey())
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, NewPublicErr(err, http.StatusUnauthorized)
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	storedIP := claims["ip"].(string)
	userID := claims["user_id"].(string)

	user, userErr := m.factory.UserRepo.GetID(c, userID)
	if userErr != nil {
		if userErr.Code == errorDb.PGErrUnexpectedError {
			return nil, NewPublicErr(err, http.StatusNotFound)
		}
		return nil, NewPublicErr(userErr.Message, http.StatusInternalServerError)
	}

	if storedIP != currentIP {
		m.log.Info("Здесь был бы вызван сервис отправки уведомления на почту", slog.String("stored_ip", storedIP), slog.String("current_ip", currentIP))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(req.RefreshToken))
	if err != nil {
		return nil, NewPublicErr(fmt.Errorf("refresh token mismatch"), http.StatusUnauthorized)
	}

	newAccessToken, err := utils.GenerateAccessToken(currentIP, userID)
	if err != nil {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	newRefreshToken, err := utils.GenerateRefreshToken(currentIP, userID)
	if err != nil {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	user.RefreshToken = newRefreshToken
	if err := m.factory.UserRepo.Update(c, user); err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
