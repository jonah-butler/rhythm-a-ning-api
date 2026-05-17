package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"rhythmapi/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var issuer = "rhythmaning-api"

const (
	JWT     = "jwt"
	REFRESH = "refresh-token"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func generateToken() (string, string, error) {
	// Generate 32 random bytes (256 bits)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", err
	}

	// Encode the random bytes to a URL-safe base64 string (without padding)
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	// Compute the SHA-256 hash of the token for storage
	tokenHash := computeHash(token)

	return token, tokenHash, nil
}

func computeHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func validateHashedPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func generateAccessClaims(user model.AuthenticateUserOutput) (*model.UserClaims, string, error) {
	fmt.Println(string(secretKey))
	claim := &model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: generateJWTRefreshExp(1), // 24 hour expiration
			Subject:   strconv.Itoa(user.UserId),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Email:    user.Email,
		Username: user.Username,
		UserId:   user.UserId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, "", err
	}

	return claim, tokenString, nil
}

func generateRefreshClaims(userClaims *model.UserClaims) (*jwt.RegisteredClaims, string, error) {
	refreshClaim := jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: generateJWTRefreshExp(14), // 2 week expiration
		Subject:   userClaims.Subject,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	signedToken, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return nil, "", err
	}

	return &refreshClaim, signedToken, nil
}

func validateToken[T jwt.Claims](token string, claimsDef T) error {
	if len(token) == 0 {
		return errors.New("no token available")
	}

	_, err := jwt.ParseWithClaims(token, claimsDef,
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		})
	if err != nil {
		return err
	}

	return nil
}

func setAuthCookies(ctx *gin.Context, jwt string, jwtExpiration int64, refreshToken string, refreshExpiration int64) {
	secure := os.Getenv("APP_ENV") == "prod"
	// set jwt cookie
	ctx.SetCookie(JWT, jwt, int(jwtExpiration), "/", "", secure, true)
	// set refresh cookie
	ctx.SetCookie(REFRESH, refreshToken, int(refreshExpiration), "/v1/user/refresh", "", secure, true)
}

func generateJWTRefreshExp(days int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add((time.Hour * 24) * time.Duration(days)).UTC())
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
