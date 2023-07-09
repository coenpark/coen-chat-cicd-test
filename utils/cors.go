package utils

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetCors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
			AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With", "custom-header", "authorization", "refreshToken", "refresh_token", "grpc-access-token", "Refreshtoken", "grpc-refresh-token"},
			AllowMethods:     []string{"POST", "GET", "PATCH", "DELETE"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
}
