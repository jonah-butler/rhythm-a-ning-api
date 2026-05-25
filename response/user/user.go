package response

import "errors"

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
