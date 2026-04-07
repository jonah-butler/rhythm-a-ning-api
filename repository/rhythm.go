package repository

import (
	"database/sql"
	"fmt"
	"rhythmapi/model"

	"github.com/gin-gonic/gin"
)

type RhythmRepository struct {
	db *sql.DB
}

func NewRhythmRepository(db *sql.DB) IRhythmRepository {
	return &RhythmRepository{db: db}
}

func (r *RhythmRepository) FindById(id int) (*model.Rhythm, error) {
	rhythm := model.Rhythm{
		Name: "hello world",
	}

	fmt.Println(id)
	return &rhythm, nil
}

func (r *RhythmRepository) GetSubdivisionTypes(ctx *gin.Context) ([]model.SubdivisionType, error) {
	var subdivisionTypes []model.SubdivisionType

	rows, err := r.db.QueryContext(ctx, GET_SUBDIVISION_TYPES)
	if err != nil {
		return subdivisionTypes, err
	}

	defer rows.Close()

	for rows.Next() {
		var subdivisionType model.SubdivisionType

		err = rows.Scan(&subdivisionType.SubdivisionId, &subdivisionType.Name)
		if err != nil {
			return subdivisionTypes, err
		}

		subdivisionTypes = append(subdivisionTypes, subdivisionType)
	}

	return subdivisionTypes, nil
}

func (r *RhythmRepository) GetRhythmLevels(ctx *gin.Context) ([]model.RhythmLevel, error) {
	var rhythmLevels []model.RhythmLevel

	rows, err := r.db.QueryContext(ctx, GET_RHYTHM_LEVELS)
	if err != nil {
		return rhythmLevels, err
	}

	for rows.Next() {
		var rhythmLevel model.RhythmLevel

		err = rows.Scan(&rhythmLevel.LevelId, &rhythmLevel.Name)
		if err != nil {
			return rhythmLevels, err
		}

		rhythmLevels = append(rhythmLevels, rhythmLevel)
	}

	return rhythmLevels, nil
}
