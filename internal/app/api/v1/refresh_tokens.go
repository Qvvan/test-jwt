// internal/app/api/v1/manager.go
package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

	currentIP := c.ClientIP()

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(utils.GetSecretKey())
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, NewPublicErr(err, http.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	storedIP, ok := claims["ip"].(string)
	if !ok {
		return nil, NewPublicErr(errors.New("invalid claim: ip"), http.StatusInternalServerError)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, NewPublicErr(errors.New("invalid claim: user_id"), http.StatusInternalServerError)
	}

	typeToken, ok := claims["type_token"].(string)
	if !ok {
		return nil, NewPublicErr(errors.New("invalid claim: type_token"), http.StatusInternalServerError)
	}

	if typeToken != "refresh_token" {
		return nil, NewPublicErr(err, http.StatusUnauthorized)
	}

	if storedIP != currentIP {
		m.log.Info("Здесь был бы вызван сервис отправки уведомления на почту", slog.String("stored_ip", storedIP), slog.String("current_ip", currentIP))
	}

	newAccessToken, err := utils.GenerateAccessToken(currentIP, userID)
	if err != nil {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	newRefreshToken, err := utils.GenerateRefreshToken(currentIP, userID)
	if err != nil {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
