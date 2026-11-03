package model

import "github.com/google/uuid"

const MaxWorkflowSections = 20

type CreateWorkflowRhythmInput struct {
	RhythmInputPoly            // all the rhythm fields, including Measures
	Section         *uuid.UUID `json:"section"` // nil if not in a section
}

// Create Payloads
type CreateSectionPayload struct {
	Color       string    `json:"color"` // hex value
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	SectionId   uuid.UUID `json:"sectionId" binding:"required"` // client generated, stored as section_id
}

// WorkflowRhythm is a rhythm plus its link data from workflow_rhythms
type WorkflowRhythm struct {
	Rhythm
	Section  *uuid.UUID `json:"section"` // nil if not in a section
	Position int        `json:"position"`
}

type CreateWorkflowPayload struct {
	Blocks      []CreateWorkflowRhythmInput `json:"blocks" binding:"dive"`
	Sections    *[]CreateSectionPayload     `json:"sections" binding:"omitempty,dive"`
	Name        string                      `json:"name" binding:"required"`
	Description *string                     `json:"description"`
}

// Create Responses
type WorkflowSection struct {
	SectionId   uuid.UUID `json:"sectionId"`
	WorkflowId  uuid.UUID `json:"workflowId"`
	Color       string    `json:"color"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	TimeStampContract
}

type Workflow struct {
	WorkflowId  uuid.UUID         `json:"workflowId"`
	UserId      uuid.UUID         `json:"userId"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	Rhythms     []WorkflowRhythm  `json:"rhythms"`
	Sections    []WorkflowSection `json:"sections"`
	TimeStampContract
}

// Update Payloads
type UpdateWorkflowRhythmInput struct {
	CreateWorkflowRhythmInput
	Id *uuid.UUID `json:"id"` // nil means a new block
}

type UpdateWorkflowInput struct {
	Name        string                      `json:"name" binding:"required"`
	Description *string                     `json:"description"`
	Blocks      []UpdateWorkflowRhythmInput `json:"blocks" binding:"dive"`
	Sections    *[]CreateSectionPayload     `json:"sections" binding:"omitempty,dive"`
}
