package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie", "Authorization", "Authtoken", "timeStamp", "language", "accept-language", "appVersion"}
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	config.AllowCredentials = true
	return cors.New(config)
}
