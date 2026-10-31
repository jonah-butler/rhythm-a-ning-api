package workflow

import (
	"errors"
)

var (
	ErrSectionCreateError = errors.New("failed to create workflow section")
	ErrWorkflowCreateErr  = errors.New("failed to create workflow")
)
