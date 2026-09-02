package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func getCloudflareSecret() string {
	return os.Getenv("CLOUDFLARE_TURNSTILE_SECRET")
}

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
