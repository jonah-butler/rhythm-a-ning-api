package route

import (
	"rhythmapi/handler"

	"github.com/gin-gonic/gin"
)

func SetupRhythmRoutes(router *gin.RouterGroup, handler *handler.RhythmHandler) {
	rhythms := router.Group("/rhythm")
	{
		rhythms.GET("/subdivisions", handler.GetSubdivisionTypes)
		rhythms.GET("/levels", handler.GetRhythmLevels)
		rhythms.GET("/:id", handler.GetRhythmById)
	}
}
