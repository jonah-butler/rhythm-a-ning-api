package repository

import (
	"context"
	"database/sql"
	"rhythmapi/model"
	"rhythmapi/request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) IWorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) CreateWorkflow(ctx *gin.Context, workflow model.CreateWorkflowPayload) {
	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return // update
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return // update
	}

	defer tx.Rollback()

	// _1 Create Workflow
	var workflowId uuid.UUID
	err = tx.QueryRowContext(
		ctx,
		CREATE_WORKFLOW,
		userId,
		workflow.Name,
		workflow.Description,
	).Scan(&workflowId)
	if err != nil {
		return // update
	}

	// _2 Create Rhythms
	rhythms := make([]model.Rhythm, len(workflow.Blocks))
	for i, block := range workflow.Blocks {
		rhythmType := block.RhythmDiscriminator.RhythmType
		if rhythmType != model.Monorhythm && rhythmType != model.Polyrhythm {
			return // update: return an "unknown rhythm type" error
		}

		r, err := insertRhythm(ctx, tx, userId, block.RhythmInputPoly, rhythmType == model.Polyrhythm)
		if err != nil {
			return // update
		}

		rhythms[i] = r
	}

	// _3 Create Sections
	var payloadSections []model.CreateSectionPayload // starts as an empty (nil) slice
	if workflow.Sections != nil {
		payloadSections = *workflow.Sections // only dereference if it was sent
	}

	_, err = createSections(ctx, tx, workflowId, payloadSections)
	if err != nil {
		return // update
	}

	// _4 Create Workflow Rhythms
	err = createWorkflowRhythms(ctx, tx, workflowId, workflow.Blocks, rhythms)
	if err != nil {
		return // update
	}

	if err := tx.Commit(); err != nil {
		return // update
	}
}

func createSections(ctx context.Context, tx *sql.Tx, workflowId uuid.UUID, sections []model.CreateSectionPayload) ([]model.WorkflowSection, error) {
	n := len(sections)
	sectionIds := make([]uuid.UUID, n)
	colors := make([]string, n)
	names := make([]string, n)
	descriptions := make([]*string, n)

	for i, s := range sections {
		sectionIds[i] = s.SectionId
		colors[i] = s.Color
		names[i] = s.Name
		descriptions[i] = s.Description
	}

	rows, err := tx.QueryContext(
		ctx,
		CREATE_SECTIONS,
		workflowId,
		pq.Array(sectionIds),
		pq.Array(colors),
		pq.Array(names),
		pq.Array(descriptions),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inserted := make([]model.WorkflowSection, 0, n)
	for rows.Next() {
		var section model.WorkflowSection
		err := rows.Scan(
			&section.SectionId,
			&section.WorkflowId,
			&section.Color,
			&section.Name,
			&section.Description,
			&section.CreatedAt,
			&section.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		inserted = append(inserted, section)
	}

	return inserted, rows.Err()
}

func createWorkflowRhythms(ctx context.Context, tx *sql.Tx, workflowId uuid.UUID, blocks []model.CreateWorkflowRhythmInput, rhythms []model.Rhythm) error {
	n := len(blocks)
	rhythmIds := make([]uuid.UUID, n)
	sectionIds := make([]*uuid.UUID, n)
	measures := make([]int, n)
	positions := make([]int, n)

	for i, block := range blocks {
		rhythmIds[i] = rhythms[i].Id
		sectionIds[i] = block.Section // nil is stored as NULL
		measures[i] = block.Measures
		positions[i] = i
	}

	_, err := tx.ExecContext(
		ctx,
		CREATE_WORKFLOW_RHYTHMS,
		workflowId,
		pq.Array(rhythmIds),
		pq.Array(sectionIds),
		pq.Array(measures),
		pq.Array(positions),
	)
	return err
}
