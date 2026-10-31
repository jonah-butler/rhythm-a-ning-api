package handler

import (
	"fmt"
	"rhythmapi/handler/workflow"
	"rhythmapi/model"

	"github.com/gin-gonic/gin"
)

func (r *WorkflowHandler) CreateWorkflow(ctx *gin.Context) {
	var workflowInput model.CreateWorkflowPayload

	if err := ctx.ShouldBindJSON(&workflowInput); err != nil {
		log := fmt.Sprintf("failed to parse workflow input: %v", err.Error())
		respondErr(ctx, log, workflow.ErrWorkflowCreateErr)
		return
	}

	newWorkflow, err := r.repo.CreateWorkflow(ctx, workflowInput)
}

func (r *WorkflowHandler) CreateSection(ctx *gin.Context) {
	var sectionPayload []model.CreateSectionPayload

	if err := ctx.ShouldBindJSON(&sectionPayload); err != nil {
		log := fmt.Sprintf("failed to parse section payload: %v", err.Error())
		respondErr(ctx, log, workflow.ErrSectionCreateError)
		return
	}

	sections, err := r.repo.CreateSection(ctx, sectionPayload)
	if err != nil {
		log := fmt.Sprintf("failed to insert section payload: %v", err.Error())
		respondErr(ctx, log, workflow.ErrSectionCreateError)
		return
	}

	respondSuccessContent(ctx, "create sections successful", gin.H{
		"sections": sections,
	})
}
