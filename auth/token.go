package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"rhythmapi/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "rhythmaning-api"

const (
	AccessTokenTTL  = 24 * time.Hour
	RefreshTokenTTL = 14 * 24 * time.Hour
)

// secretKey is read on every call rather than captured in a package-level var:
// the process environment is not populated until godotenv.Load runs inside
// main, which happens after package-level initialization.
func secretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateAccessClaims(user model.AuthenticateUserOutput) (*model.UserClaims, string, error) {
	claims := &model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: expiresIn(AccessTokenTTL),
			Subject:   user.UserId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Email:    user.Email,
		Username: user.Username,
		UserId:   user.UserId,
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey())
	if err != nil {
		return nil, "", err
	}

	return claims, signed, nil
}

func GenerateRefreshClaims(userClaims *model.UserClaims) (*jwt.RegisteredClaims, string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: expiresIn(RefreshTokenTTL),
		Subject:   userClaims.Subject,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey())
	if err != nil {
		return nil, "", err
	}

	return claims, signed, nil
}

// ValidateToken parses and verifies token into claimsDef, returning the
// populated claims. claimsDef must be a pointer to a claims struct.
func ValidateToken[T jwt.Claims](token string, claimsDef T) (T, error) {
	if len(token) == 0 {
		return claimsDef, errors.New("no token available")
	}

	parsed, err := jwt.ParseWithClaims(token, claimsDef,
		func(token *jwt.Token) (any, error) {
			return secretKey(), nil
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

// HashToken derives the value stored in user_tokens for a signed JWT. Distinct
// from handler.computeHash, which uses base64 for the emailed opaque tokens.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func expiresIn(d time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(d).UTC())
}
