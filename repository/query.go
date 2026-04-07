package repository

import (
	_ "embed"
)

//go:embed subdivisions/get_subdivision_types.sql
var GET_SUBDIVISION_TYPES string

//go:embed levels/get_rhythm_levels.sql
var GET_RHYTHM_LEVELS string
