package model

import (
	"time"

	"github.com/google/uuid"
)

type SubdivisionType struct {
	SubdivisionId int    `json:"subdivisionId"`
	Name          string `json:"name"`
}

type RhythmLevel struct {
	LevelId int    `json:"levelId"`
	Name    string `json:"name"`
}

type SubdivisionTypes string

const (
	BASE       SubdivisionTypes = "base"
	DUPLET     SubdivisionTypes = "duplet"
	TRIPLET    SubdivisionTypes = "triplet"
	QUADRUPLET SubdivisionTypes = "quadruplet"
	QUINTUPLET SubdivisionTypes = "quintuplet"
	SEXTUPLET  SubdivisionTypes = "sextuplet"
	SEPTUPLET  SubdivisionTypes = "septuplet"
	OCTUPLET   SubdivisionTypes = "octuplet"
	NONUPLET   SubdivisionTypes = "nonuplet"
	DECUPLET   SubdivisionTypes = "decuplet"
)

type RhythmType int

const (
	Monorhythm RhythmType = iota
	Polyrhythm
)

type RhythmDiscriminator struct {
	RhythmType RhythmType `json:"rhythmType" binding:"oneof=0 1"`
}

type RhythmInputMono struct {
	RhythmDiscriminator
	Bpm         int              `json:"bpm" binding:"required"`
	Beats       int              `json:"beats" binding:"required"`
	Measures    int              `json:"measures" binding:"required"`
	Subdivision SubdivisionTypes `json:"subdivision" binding:"required"`
	State       []int64          `json:"state" binding:"required"`
	Name        string           `json:"name" binding:"required"`
	Level       *string          `json:"level"`
	Description string           `json:"description"`
}

type RhythmInputPoly struct {
	RhythmInputMono
	PolyBeats       *int              `json:"polyBeats"`
	PolyState       []int64           `json:"polyState"`
	PolySubdivision *SubdivisionTypes `json:"polySubdivision"`
}

// from DB
type Rhythm struct {
	RhythmInputPoly           // contains all fields
	Id              uuid.UUID `json:"id"`
	IsPoly          bool      `json:"isPoly"`
	UserId          uuid.UUID `json:"userId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type GetRhythmsRequest struct {
	Offset int `json:"offset" binding:"required"`
}
