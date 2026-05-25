package repository

import (
	"rhythmapi/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type IRhythmRepository interface {
	GetSubdivisionTypes(ctx *gin.Context) ([]model.SubdivisionType, error)
	GetRhythmLevels(ctx *gin.Context) ([]model.RhythmLevel, error)
	GetRhythms(ctx *gin.Context, request model.GetRhythmsRequest) ([]model.Rhythm, int, error)
	GetRhythmById(ctx *gin.Context, id uuid.UUID) (model.Rhythm, error)
	DeleteRhythmById(ctx *gin.Context, id uuid.UUID) (bool, error)
	UpdateRhythmById(ctx *gin.Context, id uuid.UUID, rhythm model.RhythmInputPoly) (model.Rhythm, error)
	CreateMonoRhythm(ctx *gin.Context, rhythm model.RhythmInputMono) (model.Rhythm, error)
	CreatePolyRhythm(ctx *gin.Context, rhythm model.RhythmInputPoly) (model.Rhythm, error)
}

type IUserRepository interface {
	InsertUserRegistrationHash(ctx *gin.Context, userId string, hash string, expiration time.Time) (bool, error)
	InsertNewUser(ctx *gin.Context, user model.RegisterUserInput) (model.RegisterUserOutput, error)
	GetUserIdByUsernameAndEmail(ctx *gin.Context, username string, email string) (string, error)
	GetUserByHash(ctx *gin.Context, hash string, tokenType model.TokenType) (model.UserTokenSlim, error)
	DeleteUserRegistrationHash(ctx *gin.Context, userId string, hash string) (bool, error)
	GetUserByEmail(ctx *gin.Context, email string) (model.AuthenticateUserOutput, error)
	InsertRefreshToken(ctx *gin.Context, hashedToken string, claims *jwt.RegisteredClaims) (bool, error)
	GetUserById(ctx *gin.Context, id string) (model.AuthenticateUserOutput, error)
	InsertUserToken(ctx *gin.Context, token string, userId string, expiration time.Time, tokenType int) (bool, error)
	UpdateUserPassword(ctx *gin.Context, password string, userId string) (bool, error)
	DeleteUserToken(ctx *gin.Context, userId string, tokenType int, token string) (bool, error)
}
