package route

import (
	"database/sql"
	"rhythmapi/handler"
	"rhythmapi/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	r := gin.Default()

	rhythmRepo := repository.NewRhythmRepository(db)
	rhythmHandler := handler.NewRhythmHandler(rhythmRepo)

	v1 := r.Group("/api/v1")

	SetupRhythmRoutes(v1, rhythmHandler)

	return r
}
