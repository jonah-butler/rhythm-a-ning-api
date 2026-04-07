package repository

import (
	"rhythmapi/model"

	"github.com/gin-gonic/gin"
)

type IRhythmRepository interface {
	GetSubdivisionTypes(ctx *gin.Context) ([]model.SubdivisionType, error)
	GetRhythmLevels(ctx *gin.Context) ([]model.RhythmLevel, error)
	FindById(id int) (*model.Rhythm, error) // stubbed
}
