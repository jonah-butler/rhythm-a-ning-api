package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserBase struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username"`
}

type UserPassword struct {
	Password string `json:"password" binding:"required"`
}

type RegisterUserInput struct {
	UserBase
	UserPassword
	TurnstileToken string `json:"turnstileToken" binding:"required"`
}

type AuthenticateUserOutput struct {
	UserId   uuid.UUID `json:"userId"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type VerifyUserInput struct {
	Token string `json:"token" binding:"required"`
}

type VerifyPasswordResetInput struct {
	VerifyUserInput
	UserPassword
}

type UserTokenSlim struct {
	UserId    uuid.UUID
	ExpiresAt time.Time
}

type RegisterUserOutput struct {
	UserID        *uuid.UUID
	Email         *string
	UsernameTaken bool
	EmailTaken    bool
}

type UserLoginInput struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email    string    `json:"email"`
	Username string    `json:"username"`
	UserId   uuid.UUID `json:"userId"`
}

type TokenType int

const (
	Unknown TokenType = iota
	AccountVerification
	RefreshToken
	PasswordReset
)
