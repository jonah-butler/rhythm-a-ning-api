package response

import (
	"errors"
	"log"
	"net/http"
	response "rhythmapi/response/user"

	"github.com/gin-gonic/gin"
)

func RespondErr(c *gin.Context, logMsg string, err error) {
	log.Println(logMsg)
	c.JSON(statusFor(err), gin.H{"error": err.Error()})
}

func RespondSuccessContent(c *gin.Context, logMsg string, data gin.H) {
	log.Println(logMsg)
	c.JSON(http.StatusOK, data)
}

func ResponseSuccessNoContent(c *gin.Context, logMsg string) {
	log.Println(logMsg)
	c.Status(http.StatusNoContent)
}

func RespondBadRequest(c *gin.Context, logMsg, userMsg string) {
	log.Println(logMsg)
	c.JSON(http.StatusBadRequest, gin.H{"error": userMsg})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, response.ErrUsernameTaken), errors.Is(err, response.ErrEmailTaken):
		return http.StatusConflict
	case errors.Is(err, response.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, response.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
