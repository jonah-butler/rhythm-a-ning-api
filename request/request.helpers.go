package request

import (
	"errors"
	"rhythmapi/middlewares"
	"rhythmapi/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrNoUserInContext = errors.New("user not found in context")
)

func GetUserFromContext(ctx *gin.Context) (string, error) {
	userId, exists := ctx.Get(middlewares.USERID)
	if !exists {
		return "", ErrNoUserInContext
	}

	converted := userId.(uuid.UUID)

	return converted.String(), nil
}

func ValidatePolyRhythm(rhythm model.RhythmInputPoly) (isValid bool) {
	isValid = true

	if rhythm.RhythmType == model.Polyrhythm {
		// must contain valid poly fields
		if rhythm.PolyBeats == nil || len(rhythm.PolyState) == 0 || rhythm.PolySubdivision == nil {
			isValid = false
		}
	}

	return
}
