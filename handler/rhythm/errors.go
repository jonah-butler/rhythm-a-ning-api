package rhythm

import (
	"errors"
)

var (
	ErrRhythmCreateError = errors.New("failed to create new rhythm")
	ErrRhythmById        = errors.New("failed to find rhythm by ID")
	ErrDeleteRhythmById  = errors.New("failed to delete rhythm by id")
	ErrUpdateRhythmById  = errors.New("failed to update rhythm by ID")
	ErrGetRhythms        = errors.New("failed to get rhythms")
	ErrTagCreation       = errors.New("failed to create tags")
)
