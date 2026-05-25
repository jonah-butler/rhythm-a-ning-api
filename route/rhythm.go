package route

import (
	"rhythmapi/handler"
	m "rhythmapi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRhythmRoutes(router *gin.RouterGroup, handler *handler.RhythmHandler) {
	rhythms := router.Group("/rhythm")
	{
		rhythms.GET("/", m.AuthorizeUser(), handler.GetRhythms)
		rhythms.POST("/", m.AuthorizeUser(), handler.CreateRhtyhm)          // create new rhythm
		rhythms.GET("/:id", m.AuthorizeUser(), handler.GetRhythmById)       // get rhythm by id
		rhythms.DELETE("/:id", m.AuthorizeUser(), handler.DeleteRhythmById) // delete rhythm by id
		rhythms.PUT("/:id", m.AuthorizeUser(), handler.UpdateRhythmById)    // update rhythm by id

		rhythms.GET("/subdivisions", handler.GetSubdivisionTypes)
		rhythms.GET("/levels", handler.GetRhythmLevels)
	}
}
