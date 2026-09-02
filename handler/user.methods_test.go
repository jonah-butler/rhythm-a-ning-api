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
