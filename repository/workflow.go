package repository

import (
	"context"
	"database/sql"
	"errors"
	"rhythmapi/model"
	"rhythmapi/request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
	ErrUnknownSection    = errors.New("block references a section not in the payload")
	ErrUnknownRhythmType = errors.New("unknown rhythm type")
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) IWorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) CreateWorkflow(ctx *gin.Context, workflow model.CreateWorkflowPayload) (model.Workflow, error) {
	var created model.Workflow

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return created, err
	}

	var payloadSections []model.CreateSectionPayload
	if workflow.Sections != nil {
		payloadSections = *workflow.Sections // only dereference if it was sent
	}

	// workflow.Blocks.Section MUST exclusively contain any one of the
	// values included in workflow.Sections
	sectionIds := make(map[uuid.UUID]struct{}, len(payloadSections))
	for _, s := range payloadSections {
		sectionIds[s.SectionId] = struct{}{}
	}
	for _, block := range workflow.Blocks {
		if block.Section == nil {
			continue
		}
		if _, ok := sectionIds[*block.Section]; !ok {
			return created, ErrUnknownSection
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return created, err
	}

	defer tx.Rollback()

	// _1 Create Workflow
	err = tx.QueryRowContext(
		ctx,
		CREATE_WORKFLOW,
		userId,
		workflow.Name,
		workflow.Description,
	).Scan(
		&created.WorkflowId,
		&created.UserId,
		&created.Name,
		&created.Description,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return created, err
	}

	// _2 Create Rhythms
	rhythms := make([]model.Rhythm, len(workflow.Blocks))
	for i, block := range workflow.Blocks {
		rhythmType := block.RhythmDiscriminator.RhythmType
		if rhythmType != model.Monorhythm && rhythmType != model.Polyrhythm {
			return created, ErrUnknownRhythmType
		}

		inserted, err := insertRhythm(ctx, tx, userId, block.RhythmInputPoly, rhythmType == model.Polyrhythm)
		if err != nil {
			return created, err
		}

		rhythms[i] = inserted
	}

	// _3 Create Sections
	sections, err := createSections(ctx, tx, created.WorkflowId, payloadSections)
	if err != nil {
		return created, err
	}

	// _4 Create Workflow Rhythms
	err = createWorkflowRhythms(ctx, tx, created.WorkflowId, workflow.Blocks, rhythms)
	if err != nil {
		return created, err
	}

	if err := tx.Commit(); err != nil {
		return created, err
	}

	created.Rhythms = rhythms
	created.Sections = sections

	return created, nil
}

func (r *WorkflowRepository) UpdateWorkflow(ctx *gin.Context, workflow model.UpdateWorkflowInput, workflowId string) (model.Workflow, error) {
	var updated model.Workflow

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return updated, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return updated, err
	}

	defer tx.Rollback()

	// _1 Update Workflow
	err = tx.QueryRowContext(
		ctx,
		UPDATE_WORKFLOW,
		workflow.Name,
		workflow.Description,
		workflowId,
		userId,
	).Scan(
		&updated.WorkflowId,
		&updated.UserId,
		&updated.Name,
		&updated.Description,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return updated, err
	}

	return updated, nil
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
