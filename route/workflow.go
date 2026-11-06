package route

import (
	"rhythmapi/handler"
	m "rhythmapi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupWorkflowRoutes(router *gin.RouterGroup, handler *handler.WorkflowHandler) {
	workflows := router.Group("workflow")
	{
		workflows.GET("/", m.AuthorizeUser(), handler.GetWorkflows)             // get workflows, paged by ?offset=
		workflows.POST("/", m.AuthorizeUser(), handler.CreateWorkflow)          // create new workflow
		workflows.GET("/:id", m.AuthorizeUser(), handler.GetWorkflowById)       // get workflow by id
		workflows.PUT("/:id", m.AuthorizeUser(), handler.UpdateWorkflow)        // update workflow by id
		workflows.DELETE("/:id", m.AuthorizeUser(), handler.DeleteWorkflowById) // delete workflow by id
	}
}
