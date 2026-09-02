package middlewares

import (
	"rhythmapi/auth"
	"rhythmapi/model"
	"rhythmapi/response"

	"github.com/gin-gonic/gin"
)

const (
	USERID = "userId"
	EMAIL  = "email"
)

func AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.AccessCookie(c)
		if err != nil {
			response.RespondErr(c, "failed to retrieve user auth cookie", err)
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(token, new(model.UserClaims))
		if err != nil {
			response.RespondErr(c, "failed to validate user in authorization middleware", err)
			c.Abort()
			return
		}

		c.Set(USERID, claims.UserId)
		c.Set(EMAIL, claims.Email)
		c.Next()
	}
}
