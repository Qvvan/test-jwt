package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/qvvan/test-jwt/internal/app/api/v1"
	"github.com/qvvan/test-jwt/internal/config"
)

func InitRouters(cfg *config.Config, v1Manager *v1.Manager) *gin.Engine {
	gin.SetMode(cfg.Debug)
	router := gin.Default()

	api := router.Group("/api")
	{
		v1.RegisterV1Routes(api, v1Manager)
	}

	return router
}
