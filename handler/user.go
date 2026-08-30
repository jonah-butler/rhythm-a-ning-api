package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"rhythmapi/aws/ses"
	"rhythmapi/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const USERNAME_LENGTH_MIN = 6
const USERNAME_LENGTH_MAX = 50

func (u *UserHandler) RegisterUser(c *gin.Context) {
	var user model.RegisterUserInput

	if err := c.ShouldBindJSON(&user); err != nil {
		respondBadRequest(c, "failed to decode request: "+err.Error(), "missing required POST parameters: "+err.Error())
		return
	}

	if user.Username != "" {
		nameLength := len(user.Username)
		if nameLength < USERNAME_LENGTH_MIN || nameLength > USERNAME_LENGTH_MAX {
			respondBadRequest(c, "invalid username length", fmt.Sprintf("invalid username length - must be between %d and %d", USERNAME_LENGTH_MIN, USERNAME_LENGTH_MAX))
			return
		}
	} else {
		user.Username = user.Email
	}

	success, err := verifyTurnstileToken(c, user.TurnstileToken)
	if err != nil || !success {
		var errorMessage string
		if err.Error() == "" {
			errorMessage = "validation failed"
		} else {
			errorMessage = err.Error()
		}

		respondErr(c, "failed to validate turnstile token: "+errorMessage, ErrRegistration)
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		respondErr(c, "failed to generate password hash: "+err.Error(), ErrRegistration)
		return
	}
	user.Password = hashedPassword

	results, err := u.repo.InsertNewUser(c, user)
	if err != nil {
		respondErr(c, "failed to register user: "+err.Error(), ErrRegistration)
		return
	}

	if results.UsernameTaken {
		respondErr(c, "registration failed - username/email unavailable", ErrUsernameTaken)
		return
	}

	if results.EmailTaken {
		respondErr(c, "registration failed - email unavailable", ErrEmailTaken)
		return
	}

	token, hash, err := generateToken()
	if err != nil {
		respondErr(c, "failed to generate verification token", ErrTokenGeneration)
		return
	}

	expiration := time.Now().Add(10 * time.Minute)

	inserted, err := u.repo.InsertUserRegistrationHash(c, results.UserID.String(), hash, expiration)
	if err != nil || !inserted {
		respondErr(c, "failed to save user registration hash", ErrRegistration)
		return
	}

	emailInput, err := ses.PrepareAccountVerificationEmail(token, *results.Email)
	if err != nil {
		respondErr(c, "failed to generate registration email", ErrEmailDelivery)
		return
	}

	if err = ses.SendEmail(emailInput); err != nil {
		respondErr(c, "failed to send registration email", ErrEmailDelivery)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": inserted})
}

func (u *UserHandler) VerifyUser(c *gin.Context) {
	var userVerifyInput model.VerifyUserInput

	if err := c.ShouldBindJSON(&userVerifyInput); err != nil {
		respondBadRequest(c, "failed to decode request: "+err.Error(), "missing required POST parameters: "+err.Error())
		return
	}

	hash := computeHash(userVerifyInput.Token)

	fmt.Println("the hash: ", hash)

	user, err := u.repo.GetUserByHash(c, hash, model.AccountVerification)
	if err != nil {
		respondErr(c, "failed to find user via hash", ErrUserVerification)
		return
	}

	if user.ExpiresAt.Before(time.Now()) {
		respondErr(c, "user verification hash is expired", ErrVerificationExpired)
		return
	}

	didDelete, err := u.repo.DeleteUserRegistrationHash(c, user.UserId.String(), hash)
	if err != nil || !didDelete {
		respondErr(c, "failed to complete user registration", ErrUserVerification)
		return
	}

	responseSuccessNoContent(c, fmt.Sprintf("user registered: %d", user.UserId))
}

func (u *UserHandler) ReplayRegistration(c *gin.Context) {
	var user model.UserBase

	if err := c.ShouldBindJSON(&user); err != nil {
		respondBadRequest(c, "failed to decode request: "+err.Error(), "missing required POST parameters: "+err.Error())
		return
	}

	token, hash, err := generateToken()
	if err != nil {
		respondErr(c, "failed to generate verification token", ErrTokenGeneration)
		return
	}

	foundUser, err := u.repo.GetUserByEmail(c, user.Email)
	if err != nil {
		respondErr(c, "failed to lookup user: "+err.Error(), ErrUserNotFound)
		return
	}

	expiration := time.Now().Add(10 * time.Minute)

	updated, err := u.repo.UpdateUserRegistrationHash(c, hash, expiration, foundUser.UserId.String())
	if err != nil || !updated {
		respondErr(c, "failed to save registration hash", ErrRegistration)
		return
	}

	emailInput, err := ses.PrepareAccountVerificationEmail(token, user.Email)
	if err != nil {
		respondErr(c, "failed to generate registration email", ErrEmailDelivery)
		return
	}

	if err = ses.SendEmail(emailInput); err != nil {
		respondErr(c, "failed to send registration email", ErrEmailDelivery)
		return
	}

	responseSuccessNoContent(c, fmt.Sprintf("user registration replay successful: %d", foundUser.UserId))
}

func (u *UserHandler) AuthenticateUser(c *gin.Context) {
	var userLogin model.UserLoginInput

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		respondBadRequest(c, "failed to decode request: "+err.Error(), "missing required POST parameters: "+err.Error())
		return
	}

	user, err := u.repo.GetUserByEmail(c, userLogin.Email)
	if err != nil {
		fmt.Println(err)
		respondErr(c, "failed to lookup user by email", ErrUserNotFound)
		return
	}

	if err = validateHashedPassword(userLogin.Password, user.Password); err != nil {
		respondErr(c, "invalid password", ErrInvalidPassword)
		return
	}

	if user.AccountPending {
		respondErr(c, "account verification pending", ErrAccountPending)
		return
	}

	claims, jwt, err := generateAccessClaims(user)
	if err != nil {
		respondErr(c, "failed to generateAccessClaims", ErrLoginFailed)
		return
	}

	refreshClaims, refreshToken, err := generateRefreshClaims(claims)
	if err != nil || refreshClaims == nil {
		respondErr(c, "failed to generate refresh token", ErrLoginFailed)
		return
	}

	hashedRefreshToken := hashToken(refreshToken)

	didInsert, err := u.repo.InsertRefreshToken(c, hashedRefreshToken, refreshClaims)
	if err != nil || !didInsert {
		respondErr(c, "failed to store refresh token", ErrLoginFailed)
		return
	}

	jwtTTL := int64(time.Until(claims.ExpiresAt.Time).Seconds())
	refreshTTL := int64(time.Until(refreshClaims.ExpiresAt.Time).Seconds())

	setAuthCookies(c, jwt, jwtTTL, refreshToken, refreshTTL)

	respondSuccessContent(c, fmt.Sprintf("authenticated user: %d", user.UserId), gin.H{"success": true})
}

func (u *UserHandler) RefreshTokens(c *gin.Context) {
	token, err := c.Cookie(JWT)
	if err != nil && errors.Is(err, http.ErrNoCookie) {
		fmt.Println("no jwt cookie found")
	} else {
		jwtClaims := new(model.UserClaims)
		err = ValidateToken(token, jwtClaims)

		if err == nil {
			// JWT is still valid — no refresh needed
			respondSuccessContent(c, "token still valid", gin.H{"success": true})
			return
		}

		if !errors.Is(err, jwt.ErrTokenExpired) {
			// Invalid token for a reason other than expiry (tampered, wrong signature, etc.)
			respondErr(c, "invalid jwt: "+err.Error(), ErrUnauthorized)
			return
		}
	}

	refreshToken, err := c.Cookie(REFRESH)
	if err != nil {
		respondErr(c, "missing refresh cookie", ErrUnauthorized)
		return
	}

	refreshClaims := new(jwt.RegisteredClaims)
	if err = ValidateToken(refreshToken, refreshClaims); err != nil {
		respondErr(c, "invalid refresh token: "+err.Error(), ErrUnauthorized)
		return
	}

	hash := hashToken(refreshToken)

	user, err := u.repo.GetUserByHash(c, hash, model.RefreshToken)
	if err != nil {
		respondErr(c, "failed to get user by refresh hash: "+err.Error(), ErrUnauthorized)
		return
	}

	if user.ExpiresAt.Before(time.Now()) {
		respondErr(c, "refresh token expired", ErrUnauthorized)
		return
	}

	authenticatedUser, err := u.repo.GetUserById(c, user.UserId.String())
	if err != nil {
		respondErr(c, "failed to get user by id: "+err.Error(), ErrUnauthorized)
		return
	}

	freshClaims, freshJwt, err := generateAccessClaims(authenticatedUser)
	if err != nil {
		respondErr(c, "failed to generate access token", ErrLoginFailed)
		return
	}

	freshRefreshClaims, freshRefreshToken, err := generateRefreshClaims(freshClaims)
	if err != nil {
		respondErr(c, "failed to generate refresh token", ErrLoginFailed)
		return
	}

	hashedRefresh := hashToken(freshRefreshToken)

	didInsert, err := u.repo.InsertRefreshToken(c, hashedRefresh, freshRefreshClaims)
	if err != nil || !didInsert {
		respondErr(c, "failed to store refresh token", ErrLoginFailed)
		return
	}

	jwtTTL := int64(time.Until(freshClaims.ExpiresAt.Time).Seconds())
	refreshTTL := int64(time.Until(freshRefreshClaims.ExpiresAt.Time).Seconds())

	setAuthCookies(c, freshJwt, jwtTTL, freshRefreshToken, refreshTTL)

	respondSuccessContent(c, "tokens refreshed", gin.H{
		"userId": authenticatedUser.UserId,
		"email":  authenticatedUser.Email,
	})
}

func (u *UserHandler) ResetPassword(c *gin.Context) {
	var userInput model.UserBase // only need the email address

	if err := c.ShouldBindJSON(&userInput); err != nil {
		respondBadRequest(c, "failed to decode request: "+err.Error(), "missing required POST parameters: "+err.Error())
		return
	}

	user, err := u.repo.GetUserByEmail(c, userInput.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responseSuccessNoContent(c, "no user found for password reset")
			return
		}

		respondErr(c, "an error occurred validating user email", ErrPasswordResetFailed)
		return
	}

	token, hash, err := generateToken()
	if err != nil {
		respondErr(c, "failed to generate password reset token", ErrTokenGeneration)
		return
	}

	expiration := time.Now().Add(10 * time.Minute)

	didInsert, err := u.repo.InsertUserToken(c, hash, user.UserId.String(), expiration, int(model.PasswordReset))
	if err != nil || !didInsert {
		respondErr(c, "failed to insert user password reset token", ErrPasswordResetFailed)
		return
	}

	email, err := ses.PreparePasswordResetEmail(token, user.Email)
	if err != nil {
		respondErr(c, "failed to generate password reset email", ErrPasswordResetFailed)
		return
	}

	if err = ses.SendEmail(email); err != nil {
		// add retry eventually
		respondErr(c, "failed to send password reset email", ErrPasswordResetFailed)
		return
	}

	responseSuccessNoContent(c, "password reset email sent")
}

func (u *UserHandler) VerifyPasswordReset(c *gin.Context) {
	var passwordResetInput model.VerifyPasswordResetInput

	if err := c.ShouldBindJSON(&passwordResetInput); err != nil {
		respondBadRequest(c, "failed to decode request: "+err.Error(), "missing required POST parameters: "+err.Error())
		return
	}

	hash := computeHash(passwordResetInput.Token)

	user, err := u.repo.GetUserByHash(c, hash, model.PasswordReset)
	if err != nil {
		respondErr(c, "failed to find user via password reset hash", ErrUserVerification)
		return
	}

	if user.ExpiresAt.Before(time.Now()) {
		respondErr(c, "user verification password reset hash is expired", ErrPasswordResetFailed)
		return
	}

	hashedPassword, err := hashPassword(passwordResetInput.Password)
	if err != nil {
		respondErr(c, "failed to hash new password", ErrPasswordResetFailed)
		return
	}

	updated, err := u.repo.UpdateUserPassword(c, hashedPassword, user.UserId.String())
	if err != nil || !updated {
		respondErr(c, "failed to update user password", err)
		return
	}

	deleted, err := u.repo.DeleteUserToken(c, user.UserId.String(), int(model.PasswordReset), hash)
	if err != nil || !deleted {
		respondErr(c, "failed to delete user hash", ErrPasswordResetFailed)
		return
	}

	responseSuccessNoContent(c, fmt.Sprintf("password reset successful for: %d", user.UserId))
}

func (u *UserHandler) IdentifyUser(c *gin.Context) {
	userId, _ := c.Get("userId")
	respondSuccessContent(c, "identified user", gin.H{
		"userId": userId,
		"email":  c.GetString("email"),
	})
}

func (u *UserHandler) LogoutUser(c *gin.Context) {
	userId, exists := c.Get("userId")
	parsedId := userId.(uuid.UUID)
	if !exists {
		respondErr(c, "no user id found in context", ErrLogoutFailed)
		return
	}

	token, err := u.repo.GetUserToken(c, parsedId.String(), 2)
	if err != nil {
		respondErr(c, "no refresh token found", ErrLogoutFailed)
		return
	}

	deleted, err := u.repo.DeleteUserToken(c, parsedId.String(), 2, token)
	if err != nil || !deleted {
		respondErr(c, "failed to delete user refresh token", ErrLogoutFailed)
		return
	}

	setAuthCookies(c, "", -1, "", -1)

	responseSuccessNoContent(c, "successfully logged out user")
}
