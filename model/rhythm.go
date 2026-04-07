package model

type Rhythm struct {
	Name string `json:"name"`
}

type SubdivisionType struct {
	SubdivisionId int    `json:"subdivisionId"`
	Name          string `json:"name"`
}

type RhythmLevel struct {
	LevelId int    `json:"levelId"`
	Name    string `json:"name"`
}
