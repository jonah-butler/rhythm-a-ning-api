package repository

import (
	"database/sql"
	"fmt"
	"rhythmapi/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) InsertUserRegistrationHash(ctx *gin.Context, userId string, hash string, expiration time.Time) (bool, error) {
	var inserted bool

	err := r.db.QueryRowContext(
		ctx,
		INSERT_USER_REGISTRATION_HASH,
		userId,
		hash,
		expiration,
	).Scan(&inserted)
	if err != nil {
		return inserted, err
	}

	return inserted, nil
}

func (r *UserRepository) InsertNewUser(ctx *gin.Context, user model.RegisterUserInput) (model.RegisterUserOutput, error) {
	var registrationOutput model.RegisterUserOutput

	err := r.db.QueryRowContext(ctx, REGISTER_USER, user.Username, user.Email, user.Password).Scan(
		&registrationOutput.UserID,
		&registrationOutput.Email,
		&registrationOutput.UsernameTaken,
		&registrationOutput.EmailTaken,
	)
	if err != nil {
		return registrationOutput, err
	}

	return registrationOutput, nil
}

func (r *UserRepository) GetUserIdByUsernameAndEmail(ctx *gin.Context, username string, email string) (string, error) {
	var userId string

	err := r.db.QueryRowContext(
		ctx,
		GET_USERID_BY_USERNAME_AND_EMAIL,
		username,
		email,
	).Scan(&userId)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func (r *UserRepository) GetUserByHash(ctx *gin.Context, hash string, tokenType model.TokenType) (model.UserTokenSlim, error) {
	var user model.UserTokenSlim

	err := r.db.QueryRowContext(
		ctx,
		GET_USER_BY_HASH,
		hash,
		tokenType,
	).Scan(&user.UserId, &user.ExpiresAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUserRegistrationHash(ctx *gin.Context, userId string, hash string) (bool, error) {
	var isDeleted bool

	err := r.db.QueryRowContext(
		ctx,
		DELETE_USER_REGISTRATION_HASH,
		userId,
		hash,
	).Scan(&isDeleted)
	if err != nil {
		return isDeleted, err
	}

	return isDeleted, nil
}

func (r *UserRepository) GetUserByEmail(ctx *gin.Context, email string) (model.AuthenticateUserOutput, error) {
	var user model.AuthenticateUserOutput

	err := r.db.QueryRowContext(
		ctx,
		AUTHENTICATE_USER,
		email,
	).Scan(&user.UserId, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetUserById(ctx *gin.Context, id string) (model.AuthenticateUserOutput, error) {
	var user model.AuthenticateUserOutput

	err := r.db.QueryRowContext(
		ctx,
		GET_USER_BY_ID,
		id,
	).Scan(&user.UserId, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) InsertRefreshToken(ctx *gin.Context, hashedToken string, claims *jwt.RegisteredClaims) (bool, error) {
	var inserted bool

	err := r.db.QueryRowContext(
		ctx,
		INSERT_REFRESH_TOKEN,
		claims.Subject,
		hashedToken,
		claims.ExpiresAt.UTC(),
	).Scan(&inserted)
	if err != nil {
		fmt.Println(inserted, err)
		return inserted, err
	}

	return inserted, nil
}

func (r *UserRepository) InsertUserToken(ctx *gin.Context, token string, userId string, expiration time.Time, tokenType int) (bool, error) {
	var inserted bool

	err := r.db.QueryRowContext(
		ctx,
		INSERT_USER_TOKEN,
		token,
		userId,
		expiration,
		tokenType,
	).Scan(&inserted)
	if err != nil {
		fmt.Println(inserted, err)
		return inserted, err
	}

	return inserted, nil
}

func (r *UserRepository) UpdateUserPassword(ctx *gin.Context, password string, userId string) (bool, error) {
	var updated bool

	err := r.db.QueryRowContext(
		ctx,
		UPDATE_USER_PASSWORD,
		password,
		userId,
	).Scan(&updated)
	if err != nil {
		return updated, err
	}

	return updated, nil
}

func (r *UserRepository) DeleteUserToken(ctx *gin.Context, userId string, tokenType int, token string) (bool, error) {
	var deleted bool

	err := r.db.QueryRowContext(
		ctx,
		DELETE_USER_TOKEN,
		userId,
		tokenType,
		token,
	).Scan(&deleted)
	if err != nil {
		return deleted, err
	}

	return deleted, nil
}
