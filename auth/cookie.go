package auth

import (
	"os"

	"github.com/gin-gonic/gin"
)

const (
	AccessCookieName  = "jwt"
	RefreshCookieName = "refresh-token"

	// The refresh cookie is scoped to the one endpoint that consumes it, so it
	// is not attached to ordinary API requests.
	refreshCookiePath = "/v1/user/refresh"
)

func SetCookies(c *gin.Context, accessToken string, accessTTL int64, refreshToken string, refreshTTL int64) {
	secure := os.Getenv("APP_ENV") == "prod"
	c.SetCookie(AccessCookieName, accessToken, int(accessTTL), "/", "", secure, true)
	c.SetCookie(RefreshCookieName, refreshToken, int(refreshTTL), refreshCookiePath, "", secure, true)
}

// ClearCookies expires both auth cookies in the caller's browser.
func ClearCookies(c *gin.Context) {
	SetCookies(c, "", -1, "", -1)
}

func AccessCookie(c *gin.Context) (string, error) {
	return c.Cookie(AccessCookieName)
}

func RefreshCookie(c *gin.Context) (string, error) {
	return c.Cookie(RefreshCookieName)
}
