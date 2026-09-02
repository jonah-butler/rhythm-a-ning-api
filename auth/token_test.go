package auth

import (
	"regexp"
	"rhythmapi/model"
	"testing"

	"github.com/google/uuid"
)

func TestHashToken(t *testing.T) {
	const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk3NjQ2MjcsImlhdCI6MTc3ODU1NTAyNywiaXNzIjoicmh5dGhtYW5pbmctYXBpIiwic3ViIjoiMTQifQ.i-1oJ_4q6OZBoa38-_u1KMatAoa2TpnN1FsNP9d2ofM"

	hash := HashToken(token)

	// user_tokens.token stores this value, so the encoding must stay stable.
	if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(hash) {
		t.Fatalf("expected 64 lowercase hex chars, got %q", hash)
	}

	if again := HashToken(token); again != hash {
		t.Fatalf("HashToken is not deterministic: %q != %q", hash, again)
	}
}

func TestRefreshTokenRejectedAsAccessToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	user := model.AuthenticateUserOutput{
		UserId:   uuid.New(),
		Email:    "someone@example.com",
		Username: "someone",
	}

	accessClaims, accessToken, err := GenerateAccessClaims(user)
	if err != nil {
		t.Fatalf("GenerateAccessClaims: %v", err)
	}

	_, refreshToken, err := GenerateRefreshClaims(accessClaims)
	if err != nil {
		t.Fatalf("GenerateRefreshClaims: %v", err)
	}

	// A genuine access token must still pass, or the guard below proves nothing.
	if _, err := ValidateToken(accessToken, new(model.UserClaims)); err != nil {
		t.Fatalf("access token rejected on the access path: %v", err)
	}

	// The refresh token shares the signing key and issuer, so it parses cleanly
	// into UserClaims with every identity field zeroed. Only UserClaims.Validate
	// separates the two.
	if _, err := ValidateToken(refreshToken, new(model.UserClaims)); err == nil {
		t.Fatal("refresh token was accepted as an access token")
	}
}
