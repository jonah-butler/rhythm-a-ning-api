package handler

import (
	"errors"
	"fmt"
	"rhythmapi/handler/workflow"
	"rhythmapi/model"
	"rhythmapi/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *WorkflowHandler) CreateWorkflow(ctx *gin.Context) {
	var workflowInput model.CreateWorkflowPayload

	if err := ctx.ShouldBindJSON(&workflowInput); err != nil {
		log := fmt.Sprintf("failed to parse workflow input: %v", err.Error())
		respondBadRequest(ctx, log, workflow.ErrWorkflowCreateErr.Error())
		return
	}

	if workflowInput.Sections != nil && len(*workflowInput.Sections) > model.MaxWorkflowSections {
		log := fmt.Sprintf("too many workflow sections: %d", len(*workflowInput.Sections))
		respondBadRequest(ctx, log, fmt.Sprintf("a workflow can have at most %d sections", model.MaxWorkflowSections))
		return
	}

	newWorkflow, err := r.repo.CreateWorkflow(ctx, workflowInput)
	if err != nil {
		log := fmt.Sprintf("failed to create workflow: %v", err.Error())
		if isInvalidWorkflowInput(err) {
			respondBadRequest(ctx, log, err.Error())
			return
		}

		respondErr(ctx, log, workflow.ErrWorkflowCreateErr)
		return
	}

	respondSuccessContent(ctx, "create workflow successful", gin.H{
		"workflow": newWorkflow,
	})
}

func (r *WorkflowHandler) UpdateWorkflow(ctx *gin.Context) {
	var updateWorkflowInput model.UpdateWorkflowInput

	if err := ctx.ShouldBindBodyWithJSON(&updateWorkflowInput); err != nil {
		log := fmt.Sprintf("failed to parse update workflow input: %v", err.Error())
		respondBadRequest(ctx, log, workflow.ErrWorkflowUpdateErr.Error())
		return
	}

	if updateWorkflowInput.Sections != nil && len(*updateWorkflowInput.Sections) > model.MaxWorkflowSections {
		log := fmt.Sprintf("too many workflow sections: %d", len(*updateWorkflowInput.Sections))
		respondBadRequest(ctx, log, fmt.Sprintf("a workflow can have at most %d sections", model.MaxWorkflowSections))
		return
	}

	workflowId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log := fmt.Sprintf("invalid workflow id: %v", err.Error())
		respondBadRequest(ctx, log, "invalid workflow id")
		return
	}

	updatedWorkflow, err := r.repo.UpdateWorkflow(ctx, updateWorkflowInput, workflowId)
	if err != nil {
		log := fmt.Sprintf("failed to update workflow: %v", err.Error())
		if errors.Is(err, repository.ErrWorkflowNotFound) {
			respondErr(ctx, log, err)
			return
		}
		if isInvalidWorkflowInput(err) {
			respondBadRequest(ctx, log, err.Error())
			return
		}

		respondErr(ctx, log, workflow.ErrWorkflowUpdateErr)
		return
	}

	respondSuccessContent(ctx, "update workflow successful", gin.H{
		"workflow": updatedWorkflow,
	})
}

func isInvalidWorkflowInput(err error) bool {
	return errors.Is(err, repository.ErrUnknownSection) ||
		errors.Is(err, repository.ErrUnknownRhythmType) ||
		errors.Is(err, repository.ErrUnknownRhythm) ||
		errors.Is(err, repository.ErrDuplicateSection) ||
		errors.Is(err, repository.ErrDuplicateRhythm) ||
		errors.Is(err, repository.ErrSectionConflict)
}

func (r *WorkflowHandler) GetWorkflowById(ctx *gin.Context) {
	workflowId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log := fmt.Sprintf("invalid workflow id: %v", err.Error())
		respondBadRequest(ctx, log, "invalid workflow id")
		return
	}

	foundWorkflow, err := r.repo.GetWorkflowById(ctx, workflowId)
	if err != nil {
		log := fmt.Sprintf("failed to get workflow %v: %v", workflowId.String(), err.Error())
		if errors.Is(err, repository.ErrWorkflowNotFound) {
			respondErr(ctx, log, err)
			return
		}

		respondErr(ctx, log, workflow.ErrWorkflowGetErr)
		return
	}

	respondSuccessContent(ctx, "get workflow successful", gin.H{
		"workflow": foundWorkflow,
	})
}

func (r *WorkflowHandler) GetWorkflows(ctx *gin.Context) {
	var workflowsRequest model.GetWorkflowsRequest

	offset := ctx.Query("offset")
	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil || parsedOffset < 0 {
			log := fmt.Sprintf("invalid offset: %q", offset)
			respondBadRequest(ctx, log, "invalid offset")
			return
		}

		workflowsRequest.Offset = parsedOffset
	}

	workflows, total, err := r.repo.GetWorkflows(ctx, workflowsRequest)
	if err != nil {
		log := fmt.Sprintf("failed to get workflows at offset %d: %v", workflowsRequest.Offset, err.Error())
		respondErr(ctx, log, workflow.ErrWorkflowsGetErr)
		return
	}

	respondSuccessContent(ctx, "successfully got workflows", gin.H{"workflows": workflows, "total": total})
}

func (r *WorkflowHandler) DeleteWorkflowById(ctx *gin.Context) {
	workflowId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log := fmt.Sprintf("invalid workflow id: %v", err.Error())
		respondBadRequest(ctx, log, "invalid workflow id")
		return
	}

	err = r.repo.DeleteWorkflowById(ctx, workflowId)
	if err != nil {
		log := fmt.Sprintf("failed to delete workflow %v: %v", workflowId.String(), err.Error())
		if errors.Is(err, repository.ErrWorkflowNotFound) {
			respondErr(ctx, log, err)
			return
		}

		respondErr(ctx, log, workflow.ErrWorkflowDeleteErr)
		return
	}

	responseSuccessNoContent(ctx, fmt.Sprintf("successfully deleted workflow %v", workflowId.String()))
}
