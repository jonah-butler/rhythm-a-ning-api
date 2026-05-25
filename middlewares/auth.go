package middlewares

import (
	"rhythmapi/model"
	"rhythmapi/response"

	"github.com/gin-gonic/gin"
)

const (
	USERID = "userId"
)

func AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := GetAuthCookie(c)
		if err != nil {
			response.RespondErr(c, "failed to retrieve user auth cookie", err)
			c.Abort()
			return
		}

		jwtClaims := new(model.UserClaims)
		claims, err := ValidateToken(jwt, jwtClaims)
		if err != nil {
			response.RespondErr(c, "failed to validate user in authorization middleware", err)
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Next()
	}
}
