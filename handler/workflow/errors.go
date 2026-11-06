package workflow

import (
	"errors"
)

var (
	ErrWorkflowCreateErr = errors.New("failed to create workflow")
	ErrWorkflowUpdateErr = errors.New("failed to update workflow")
	ErrWorkflowGetErr    = errors.New("failed to get workflow")
	ErrWorkflowsGetErr   = errors.New("failed to get workflows")
	ErrWorkflowDeleteErr = errors.New("failed to delete workflow")
)
