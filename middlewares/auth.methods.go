package middlewares

import (
	"errors"
	"fmt"
	"os"
	"rhythmapi/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var issuer = "rhythmaning-api"

const (
	JWT     = "jwt"
	REFRESH = "refresh-token"
)

func GenerateAccessClaims(user model.AuthenticateUserOutput) (*model.UserClaims, string, error) {
	fmt.Println(string(secretKey))
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
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, "", err
	}

	return claim, tokenString, nil
}

func GenerateRefreshClaims(userClaims *model.UserClaims) (*jwt.RegisteredClaims, string, error) {
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

func ValidateToken[T jwt.Claims](token string, claimsDef T) (T, error) {
	if len(token) == 0 {
		return claimsDef, errors.New("no token available")
	}

	parsed, err := jwt.ParseWithClaims(token, claimsDef,
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		})
	if err != nil {
		return claimsDef, err
	}

	claims, ok := parsed.Claims.(T)
	if !ok {
		return claimsDef, errors.New("invalid claims type")
	}

	return claims, nil
}

func SetAuthCookies(ctx *gin.Context, jwt string, jwtExpiration int64, refreshToken string, refreshExpiration int64) {
	secure := os.Getenv("APP_ENV") == "prod"
	// set jwt cookie
	ctx.SetCookie(JWT, jwt, int(jwtExpiration), "/", "", secure, true)
	// set refresh cookie
	ctx.SetCookie(REFRESH, refreshToken, int(refreshExpiration), "/v1/user/refresh", "", secure, true)
}

func GetAuthCookie(ctx *gin.Context) (jwt string, err error) {
	jwt, err = ctx.Cookie(JWT)
	if err != nil {
		return
	}

	return
}

func generateJWTRefreshExp(days int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add((time.Hour * 24) * time.Duration(days)).UTC())
}
