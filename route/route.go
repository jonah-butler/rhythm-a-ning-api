package route

import (
	"database/sql"
	"rhythmapi/handler"
	"rhythmapi/middlewares"
	"rhythmapi/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORS())

	// repo initialization
	rhythmRepo := repository.NewRhythmRepository(db)
	userRepo := repository.NewUserRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)

	// handler initialization
	rhythmHandler := handler.NewRhythmHandler(rhythmRepo)
	userHandler := handler.NewUserHandler(userRepo)
	workflowHandler := handler.NewWorkflowHandler(workflowRepo)

	// leading api groups
	v1 := r.Group("/v1")

	// route initialization
	SetupRhythmRoutes(v1, rhythmHandler)
	SetupUserRoutes(v1, userHandler)
	SetupWorkflowRoutes(v1, workflowHandler)

	return r
}
