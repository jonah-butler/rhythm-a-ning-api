package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

var validSubdivisions = map[SubdivisionTypes]struct{}{
	BASE: {}, DUPLET: {}, TRIPLET: {}, QUADRUPLET: {}, QUINTUPLET: {},
	SEXTUPLET: {}, SEPTUPLET: {}, OCTUPLET: {}, NONUPLET: {}, DECUPLET: {},
}

func (s SubdivisionTypes) IsValid() bool {
	_, ok := validSubdivisions[s]
	return ok
}

// Subdivisions holds one subdivision per beat, stored as TEXT[]
type Subdivisions []SubdivisionTypes

func (s Subdivisions) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	arr := make(pq.StringArray, len(s))
	for i, sub := range s {
		arr[i] = string(sub)
	}
	return arr.Value()
}

func (s *Subdivisions) Scan(value any) error {
	var arr pq.StringArray
	if err := arr.Scan(value); err != nil {
		return err
	}

	if arr == nil {
		*s = nil
		return nil
	}

	subs := make(Subdivisions, len(arr))
	for i, sub := range arr {
		subs[i] = SubdivisionTypes(sub)
	}
	*s = subs
	return nil
}

type RhythmType int

const (
	Monorhythm RhythmType = iota
	Polyrhythm
)

type RhythmDiscriminator struct {
	RhythmType RhythmType `json:"rhythmType" binding:"oneof=0 1"`
}

type JSONBSounds [][]int16

func (j JSONBSounds) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONBSounds) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONBSounds")
	}
	return json.Unmarshal(bytes, j)
}

type RhythmInputMono struct {
	RhythmDiscriminator
	Bpm          int              `json:"bpm" binding:"required"`
	Beats        int              `json:"beats" binding:"required"`
	Measures     int              `json:"measures" binding:"required"`
	Sounds       JSONBSounds      `json:"sounds"`
	Subdivision  SubdivisionTypes `json:"subdivision" binding:"required"`
	Subdivisions Subdivisions     `json:"subdivisions" binding:"required"`
	State        []int64          `json:"state" binding:"required"`
	Name         string           `json:"name" binding:"required"`
	Level        *string          `json:"level"`
	Description  string           `json:"description"`
}

type RhythmInputPoly struct {
	RhythmInputMono
	PolyBeats        *int              `json:"polyBeats"`
	PolyState        []int64           `json:"polyState"`
	PolySounds       JSONBSounds       `json:"polySounds"`
	PolySubdivision  *SubdivisionTypes `json:"polySubdivision"`
	PolySubdivisions Subdivisions      `json:"polySubdivisions"`
}

type CreateTagInput struct {
	Tags []string `json:"tags"`
}

// from DB
type Rhythm struct {
	RhythmInputPoly           // contains all fields
	Id              uuid.UUID `json:"id"`
	IsPoly          bool      `json:"isPoly"`
	UserId          uuid.UUID `json:"userId"`
	TimeStampContract
}

type GetRhythmsRequest struct {
	Offset int `json:"offset" binding:"required"`
}
