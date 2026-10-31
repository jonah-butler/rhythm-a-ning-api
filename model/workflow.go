package model

import "github.com/google/uuid"

type CreateWorkflowRhythmInput struct {
	RhythmInputPoly            // all the rhythm fields, including Measures
	Section         *uuid.UUID `json:"section"` // nil if not in a section
}

// Create Payloads
type CreateSectionPayload struct {
	WorkflowId  uuid.UUID `json:"workflowId"`
	Color       string    `json:"color"` // hex value
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	SectionId   uuid.UUID `json:"sectionId"` // local state uuid
}

type WorkflowRhythm struct {
	Section *uuid.UUID `json:"section"`
	Name    *string    `json:"name"`
	Rhythm
}

type CreateWorkflowPayload struct {
	Blocks      []CreateWorkflowRhythmInput `json:"blocks"`
	Sections    *[]CreateSectionPayload     `json:"sections"`
	Name        string                      `json:"name"`
	Description *string                     `json:"description"`
}

// Create Responses
type WorkflowSection struct {
	SectionId uuid.UUID `json:"sectionId"`
	CreateSectionPayload
	TimeStampContract
}
