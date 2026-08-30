package route

import (
	"rhythmapi/handler"
	m "rhythmapi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, handler *handler.UserHandler) {
	user := router.Group("/user")
	{
		user.POST("/", handler.RegisterUser)                             // register a new user
		user.POST("/verify", handler.VerifyUser)                         // verify new user account
		user.POST("/replay-registration", handler.ReplayRegistration)    // replay user verification step
		user.POST("/login", handler.AuthenticateUser)                    // login user
		user.POST("/refresh", handler.RefreshTokens)                     // generates JWT & Refresh Token for user
		user.POST("/reset-password", handler.ResetPassword)              // generates reset token and emails
		user.POST("/verify-password-reset", handler.VerifyPasswordReset) // verify and complete the password reset flow
		user.GET("/me", m.AuthorizeUser(), handler.IdentifyUser)         // validates token and returns user id and email
		user.GET("/logout", m.AuthorizeUser(), handler.LogoutUser)
	}
}
