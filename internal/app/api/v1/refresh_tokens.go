// internal/app/api/v1/manager.go
package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qvvan/test-jwt/internal/app/utils"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
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

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, NewPublicErr(err, http.StatusBadRequest)
	}

	storedIP, userID, currentHashToken, err := utils.DecodeUserData(req.RefreshToken)
	if err != nil {
		return nil, NewPublicErr(err, http.StatusBadRequest)
	}

	user, userErr := m.factory.UserRepo.GetID(userID)
	if userErr != nil {
		if userErr == sql.ErrNoRows {
			return nil, NewPublicErr(err, http.StatusNotFound)
		}
		return nil, err
	}

	if user.RefreshToken != currentHashToken {
		return nil, NewPublicErr(fmt.Errorf("invalid refresh token"), http.StatusUnauthorized)
	}

	currentIP := c.ClientIP()

	if storedIP != currentIP {
		m.log.Info("Здесь был бы вызван сервис отправки уведомления на почту", slog.String("stored_ip", storedIP), slog.String("current_ip", currentIP))
	}

	newAccessToken, err := utils.GenerateAccessToken(currentIP, userID)
	if err != nil {
		return nil, err
	}

	hashToken, newRefreshToken := utils.GenerateRefreshToken(currentIP, userID)

	user.RefreshToken = hashToken

	if err := m.factory.UserRepo.Update(user); err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
