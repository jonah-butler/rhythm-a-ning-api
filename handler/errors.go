package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUsernameTaken       = errors.New("username unavailable")
	ErrEmailTaken          = errors.New("email unavailable")
	ErrUserNotFound        = errors.New("user not found")
	ErrRegistration        = errors.New("registration failed")
	ErrTokenGeneration     = errors.New("token generation failed")
	ErrEmailDelivery       = errors.New("failed to send verification email")
	ErrUserVerification    = errors.New("failed to verify user")
	ErrVerificationExpired = errors.New("user verification expired")
	ErrInvalidPassword     = errors.New("failed to validate password")
	ErrLoginFailed         = errors.New("failed to login user")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrPasswordResetFailed = errors.New("password reset failed")
)

func respondErr(c *gin.Context, logMsg string, err error) {
	log.Println(logMsg)
	c.JSON(statusFor(err), gin.H{"error": err.Error()})
}

func respondSuccessContent(c *gin.Context, logMsg string, data gin.H) {
	log.Println(logMsg)
	c.JSON(http.StatusOK, data)
}

func responseSuccessNoContent(c *gin.Context, logMsg string) {
	log.Println(logMsg)
	c.Status(http.StatusNoContent)
}
func respondBadRequest(c *gin.Context, logMsg, userMsg string) {
	log.Println(logMsg)
	c.JSON(http.StatusBadRequest, gin.H{"error": userMsg})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, ErrUsernameTaken), errors.Is(err, ErrEmailTaken):
		return http.StatusConflict
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
