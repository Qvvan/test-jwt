package v1

import "github.com/gin-gonic/gin"

func RegisterV1Routes(apiGroup *gin.RouterGroup, manager *Manager) {
	v1Group := apiGroup.Group("/v1")

	auth := v1Group.Group("/auth")
	{
		auth.POST("/tokens", manager.GetTokens)
		auth.POST("/tokens/refresh", manager.RefreshTokens)
	}
}
