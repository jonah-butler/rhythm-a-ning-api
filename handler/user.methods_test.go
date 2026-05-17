package handler

import (
	"fmt"
	"testing"
)

func TestComputeTokenGeneration(t *testing.T) {
	token, hash, err := generateToken()
	if err != nil {
		fmt.Println(nil)
	}

	fmt.Println(token, hash)
}

func TestComputeHash(t *testing.T) {
	token := "kWBEpCRLevL6oRiDWhtGlSXrIlyfJQ5T310OMr4Ql7I"

	hash := computeHash(token)
	fmt.Println(hash)
}

func TestHashToken(t *testing.T) {
	hashedtoken := hashToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk3NjQ2MjcsImlhdCI6MTc3ODU1NTAyNywiaXNzIjoicmh5dGhtYW5pbmctYXBpIiwic3ViIjoiMTQifQ.i-1oJ_4q6OZBoa38-_u1KMatAoa2TpnN1FsNP9d2ofM")
	fmt.Println(hashedtoken)
}
