package repository

import (
	"rhythmapi/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type IRhythmRepository interface {
	GetSubdivisionTypes(ctx *gin.Context) ([]model.SubdivisionType, error)
	GetRhythmLevels(ctx *gin.Context) ([]model.RhythmLevel, error)
	FindById(id int) (*model.Rhythm, error) // stubbed
}

type IUserRepository interface {
	InsertUserRegistrationHash(ctx *gin.Context, userId int, hash string, expiration time.Time) (bool, error)
	InsertNewUser(ctx *gin.Context, user model.RegisterUserInput) (model.RegisterUserOutput, error)
	GetUserIdByUsernameAndEmail(ctx *gin.Context, username string, email string) (int, error)
	GetUserByHash(ctx *gin.Context, hash string, tokenType model.TokenType) (model.UserTokenSlim, error)
	DeleteUserRegistrationHash(ctx *gin.Context, userId int, hash string) (bool, error)
	GetUserByEmail(ctx *gin.Context, email string) (model.AuthenticateUserOutput, error)
	InsertRefreshToken(ctx *gin.Context, hashedToken string, claims *jwt.RegisteredClaims) (bool, error)
	GetUserById(ctx *gin.Context, id int) (model.AuthenticateUserOutput, error)
	InsertUserToken(ctx *gin.Context, token string, userId int, expiration time.Time, tokenType int) (bool, error)
	UpdateUserPassword(ctx *gin.Context, password string, userId int) (bool, error)
	DeleteUserToken(ctx *gin.Context, userId int, tokenType int, token string) (bool, error)
}
