package request

import (
	"errors"
	"fmt"
	"rhythmapi/middlewares"
	"rhythmapi/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrNoUserInContext = errors.New("user not found in context")
)

func GetUserFromContext(ctx *gin.Context) (*uuid.UUID, error) {
	userId, exists := ctx.Get(middlewares.USERID)
	if !exists {
		return nil, ErrNoUserInContext
	}

	converted := userId.(uuid.UUID)
	return &converted, nil
}

func ValidatePolyRhythm(rhythm model.RhythmInputPoly) (isValid bool) {
	isValid = true

	if rhythm.RhythmType == model.Polyrhythm {
		// must contain valid poly fields
		if rhythm.PolyBeats == nil || len(rhythm.PolyState) == 0 || rhythm.PolySubdivision == nil || rhythm.PolySubdivisions == nil {
			isValid = false
		}
	}

	return
}

// ValidateSubdivisions ensures there is exactly one known subdivision per beat
func ValidateSubdivisions(beats int, subdivisions model.Subdivisions) error {
	if len(subdivisions) != beats {
		return fmt.Errorf("expected %d subdivisions (one per beat), got %d", beats, len(subdivisions))
	}

	for i, sub := range subdivisions {
		if !sub.IsValid() {
			return fmt.Errorf("invalid subdivision %q at beat %d", sub, i)
		}
	}

	return nil
}

func ValidateMonoRhythm(rhythm model.RhythmInputMono) error {
	if err := ValidateSubdivisions(rhythm.Beats, rhythm.Subdivisions); err != nil {
		return fmt.Errorf("subdivisions: %w", err)
	}

	return nil
}

func ValidateRhythm(rhythm model.RhythmInputPoly) error {
	if err := ValidateMonoRhythm(rhythm.RhythmInputMono); err != nil {
		return err
	}

	if rhythm.RhythmType != model.Polyrhythm {
		return nil
	}

	if !ValidatePolyRhythm(rhythm) {
		return errors.New("poly rhythms must contain associated poly fields")
	}

	if err := ValidateSubdivisions(*rhythm.PolyBeats, rhythm.PolySubdivisions); err != nil {
		return fmt.Errorf("polySubdivisions: %w", err)
	}

	return nil
}
