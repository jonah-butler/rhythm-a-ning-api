package route

import (
	"rhythmapi/handler"
	m "rhythmapi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupWorkflowRoutes(router *gin.RouterGroup, handler *handler.WorkflowHandler) {
	workflows := router.Group("workflow")
	{
		workflows.POST("/", m.AuthorizeUser(), handler.CreateWorkflow)
		workflows.PUT("/:id", m.AuthorizeUser(), handler.UpdateWorkflow)
	}
}
