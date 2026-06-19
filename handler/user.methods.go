package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"rhythmapi/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func getSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func getCloudflareSecret() string {
	return os.Getenv("CLOUDFLARE_TURNSTILE_SECRET")
}

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
	claim := &model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: generateJWTRefreshExp(1), // 24 hour expiration
			Subject:   user.UserId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Email:    user.Email,
		Username: user.Username,
		UserId:   user.UserId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(getSecretKey())
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
	signedToken, err := refreshToken.SignedString(getSecretKey())
	if err != nil {
		return nil, "", err
	}

	return &refreshClaim, signedToken, nil
}

func ValidateToken[T jwt.Claims](token string, claimsDef T) error {
	if len(token) == 0 {
		return errors.New("no token available")
	}

	_, err := jwt.ParseWithClaims(token, claimsDef,
		func(token *jwt.Token) (any, error) {
			return getSecretKey(), nil
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

func GetAuthCookies(ctx *gin.Context) (jwt, refresh string, err error) {
	jwt, err = ctx.Cookie(JWT)
	if err != nil {
		return
	}

	refresh, err = ctx.Cookie(REFRESH)
	return
}

func generateJWTRefreshExp(days int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add((time.Hour * 24) * time.Duration(days)).UTC())
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

type TurnstileResponse struct {
	Success bool `json:"success"`
}

func verifyTurnstileToken(ctx *gin.Context, turnstileToken string) (bool, error) {
	payload := map[string]any{
		"secret":   getCloudflareSecret(),
		"response": turnstileToken,
	}

	fmt.Println(payload)

	body, _ := json.Marshal(payload)

	endpoint := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	req, err := http.NewRequestWithContext(ctx.Request.Context(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		fmt.Println("error building request", err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error response", err)
		return false, err
	}

	defer resp.Body.Close()

	var result TurnstileResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("error decoding response", err)
		return false, err
	}

	return result.Success, nil
}
