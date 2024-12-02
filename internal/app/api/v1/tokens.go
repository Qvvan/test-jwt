package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qvvan/test-jwt/internal/app/utils"
	errorDb "github.com/qvvan/test-jwt/pkg/client/postgresql/utils"
)

type TokensRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (m *Manager) GetTokens(c *gin.Context) {
	resp, err := m.GetTokensService(c)

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

func (m *Manager) GetTokensService(c *gin.Context) (*TokensResponse, error) {
	var req TokensRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, NewPublicErr(err, http.StatusBadRequest)
	}

	user, err := m.factory.UserRepo.GetID(c, req.UserID)
	if err != nil {
		if err.Code == errorDb.PGErrUnexpectedError {
			return nil, NewPublicErr(err, http.StatusNotFound)
		}
		return nil, NewPublicErr(err.Message, http.StatusInternalServerError)
	}

	currentIP := c.ClientIP()
	accessToken, errToken := utils.GenerateAccessToken(currentIP, user.UserID)
	if errToken != nil {
		return nil, NewPublicErr(errToken, http.StatusInternalServerError)
	}

	refreshToken, errToken := utils.GenerateRefreshToken(currentIP, user.UserID)
	if errToken != nil {
		return nil, NewPublicErr(errToken, http.StatusInternalServerError)
	}

	user.RefreshToken = refreshToken

	if err := m.factory.UserRepo.Update(c, user); err != nil {
		return nil, NewPublicErr(err, http.StatusInternalServerError)
	}

	return &TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
