package workflow

import (
	"errors"
)

var (
	ErrWorkflowCreateErr = errors.New("failed to create workflow")
	ErrWorkflowUpdateErr = errors.New("failed to update workflow")
)
