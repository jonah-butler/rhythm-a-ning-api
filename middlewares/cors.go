package middlewares

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	allowedOrigins := parseOrigins(os.Getenv("ALLOWED_ORIGINS"))

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{}
	}
	origins := strings.Split(raw, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}
	return origins
}
