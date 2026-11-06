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
	ErrUnknownRhythm     = errors.New("block references a rhythm not in this workflow")
	ErrDuplicateSection  = errors.New("section id appears more than once")
	ErrDuplicateRhythm   = errors.New("rhythm id appears more than once")
	ErrSectionConflict   = errors.New("section id belongs to another workflow")
	ErrWorkflowNotFound  = errors.New("workflow not found")
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
	blockSections := make([]*uuid.UUID, len(workflow.Blocks))
	for i, block := range workflow.Blocks {
		blockSections[i] = block.Section
	}
	if err := validateSections(payloadSections, blockSections); err != nil {
		return created, err
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

		inserted.Measures = block.Measures // measures lives on workflow_rhythms, not rhythms
		rhythms[i] = inserted
	}

	// _3 Create Sections
	sections, err := saveSections(ctx, tx, CREATE_SECTIONS, created.WorkflowId, payloadSections)
	if err != nil {
		return created, err
	}

	// _4 Create Workflow Rhythms
	err = saveWorkflowRhythms(ctx, tx, CREATE_WORKFLOW_RHYTHMS, created.WorkflowId, workflow.Blocks, rhythms)
	if err != nil {
		return created, err
	}

	if err := tx.Commit(); err != nil {
		return created, err
	}

	created.Rhythms = toWorkflowRhythms(workflow.Blocks, rhythms)
	created.Sections = sections

	return created, nil
}

func (r *WorkflowRepository) UpdateWorkflow(ctx *gin.Context, workflow model.UpdateWorkflowInput, workflowId uuid.UUID) (model.Workflow, error) {
	var updated model.Workflow

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return updated, err
	}

	var payloadSections []model.CreateSectionPayload
	if workflow.Sections != nil {
		payloadSections = *workflow.Sections // only dereference if it was sent
	}

	blockSections := make([]*uuid.UUID, len(workflow.Blocks))
	blocks := make([]model.CreateWorkflowRhythmInput, len(workflow.Blocks))
	blockIds := make(map[uuid.UUID]struct{}, len(workflow.Blocks))
	for i, block := range workflow.Blocks {
		blockSections[i] = block.Section
		blocks[i] = block.CreateWorkflowRhythmInput

		if block.Id == nil {
			continue
		}
		if _, ok := blockIds[*block.Id]; ok {
			return updated, ErrDuplicateRhythm
		}
		blockIds[*block.Id] = struct{}{}
	}
	if err := validateSections(payloadSections, blockSections); err != nil {
		return updated, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return updated, err
	}

	defer tx.Rollback()

	// _1 Update Workflow: checks ownership and locks the row until commit
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
	if errors.Is(err, sql.ErrNoRows) {
		return updated, ErrWorkflowNotFound
	}
	if err != nil {
		return updated, err
	}

	// _2 Read the rhythms currently in this workflow (before the links are rebuilt)
	existingIds, err := getWorkflowRhythmIds(ctx, tx, workflowId)
	if err != nil {
		return updated, err
	}

	// _3 Sync Sections: upsert those sent, delete those left out
	sections, err := saveSections(ctx, tx, UPSERT_SECTIONS, workflowId, payloadSections)
	if err != nil {
		return updated, err
	}
	if len(sections) != len(payloadSections) {
		return updated, ErrSectionConflict
	}

	keptSectionIds := make([]uuid.UUID, len(payloadSections)) // non-nil so an empty list is {} not NULL
	for i, s := range payloadSections {
		keptSectionIds[i] = s.SectionId
	}
	_, err = tx.ExecContext(ctx, DELETE_REMOVED_SECTIONS, workflowId, pq.Array(keptSectionIds))
	if err != nil {
		return updated, err
	}

	// _4 Sync Rhythms: update blocks with an id, insert blocks without one
	rhythms := make([]model.Rhythm, len(workflow.Blocks))
	for i, block := range workflow.Blocks {
		rhythmType := block.RhythmDiscriminator.RhythmType
		if rhythmType != model.Monorhythm && rhythmType != model.Polyrhythm {
			return updated, ErrUnknownRhythmType
		}

		if block.Id == nil {
			inserted, err := insertRhythm(ctx, tx, userId, block.RhythmInputPoly, rhythmType == model.Polyrhythm)
			if err != nil {
				return updated, err
			}

			inserted.Measures = block.Measures
			rhythms[i] = inserted
			continue
		}

		if _, ok := existingIds[*block.Id]; !ok {
			return updated, ErrUnknownRhythm
		}

		saved, err := updateRhythm(ctx, tx, *block.Id, userId, block.RhythmInputPoly)
		if err != nil {
			return updated, err
		}

		saved.Measures = block.Measures
		rhythms[i] = saved
	}

	// _5 Sync Workflow Rhythms: delete links for removed blocks, then upsert the rest
	// in place so each block keeps its workflow_rhythm_id
	keptRhythmIds := make([]uuid.UUID, len(rhythms)) // non-nil so an empty list is {} not NULL
	for i, rhythm := range rhythms {
		keptRhythmIds[i] = rhythm.Id
	}
	_, err = tx.ExecContext(ctx, DELETE_REMOVED_WORKFLOW_RHYTHMS, workflowId, pq.Array(keptRhythmIds))
	if err != nil {
		return updated, err
	}

	err = saveWorkflowRhythms(ctx, tx, UPSERT_WORKFLOW_RHYTHMS, workflowId, blocks, rhythms)
	if err != nil {
		return updated, err
	}

	// _6 Delete rhythms that were removed from the workflow
	oldIds := make([]uuid.UUID, 0, len(existingIds))
	for id := range existingIds {
		oldIds = append(oldIds, id)
	}
	_, err = tx.ExecContext(ctx, DELETE_ORPHAN_RHYTHMS, pq.Array(oldIds), userId)
	if err != nil {
		return updated, err
	}

	if err := tx.Commit(); err != nil {
		return updated, err
	}

	updated.Rhythms = toWorkflowRhythms(blocks, rhythms)
	updated.Sections = sections

	return updated, nil
}

func (r *WorkflowRepository) GetWorkflowById(ctx *gin.Context, workflowId uuid.UUID) (model.Workflow, error) {
	var workflow model.Workflow

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return workflow, err
	}

	// one read-only transaction so the workflow, rhythms and sections come
	// from the same snapshot even if an update commits in between
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return workflow, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(
		ctx,
		GET_WORKFLOW_BY_ID,
		workflowId,
		userId,
	).Scan(
		&workflow.WorkflowId,
		&workflow.UserId,
		&workflow.Name,
		&workflow.Description,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return workflow, ErrWorkflowNotFound
	}
	if err != nil {
		return workflow, err
	}

	workflow.Rhythms, err = getWorkflowRhythms(ctx, tx, workflowId)
	if err != nil {
		return workflow, err
	}

	workflow.Sections, err = getWorkflowSections(ctx, tx, workflowId)
	if err != nil {
		return workflow, err
	}

	return workflow, tx.Commit()
}

func (r *WorkflowRepository) GetWorkflows(ctx *gin.Context, workflowsRequest model.GetWorkflowsRequest) ([]model.WorkflowListItem, int, error) {
	workflows := make([]model.WorkflowListItem, 0)

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return workflows, 0, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		GET_WORKFLOWS,
		userId,
		workflowsRequest.Offset,
	)
	if err != nil {
		return workflows, 0, err
	}

	defer rows.Close()

	var total int
	for rows.Next() {
		var workflow model.WorkflowListItem

		err = rows.Scan(
			&workflow.WorkflowId,
			&workflow.UserId,
			&workflow.Name,
			&workflow.Description,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
			&workflow.BlockCount,
			&workflow.TotalMeasures,
			&workflow.AverageTempo,
			&workflow.TotalMinutes,
			&total,
		)
		if err != nil {
			return workflows, total, err
		}

		workflows = append(workflows, workflow)
	}

	return workflows, total, rows.Err()
}

func (r *WorkflowRepository) DeleteWorkflowById(ctx *gin.Context, workflowId uuid.UUID) error {
	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// _1 Check ownership and lock the workflow so an update can't add rhythms mid-delete
	var lockedId uuid.UUID
	err = tx.QueryRowContext(ctx, LOCK_WORKFLOW, workflowId, userId).Scan(&lockedId)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrWorkflowNotFound
	}
	if err != nil {
		return err
	}

	// _2 Delete the workflow's rhythms; their workflow_rhythms rows cascade
	_, err = tx.ExecContext(ctx, DELETE_WORKFLOW_BLOCK_RHYTHMS, workflowId, userId)
	if err != nil {
		return err
	}

	// _3 Delete the workflow; its sections cascade
	_, err = tx.ExecContext(ctx, DELETE_WORKFLOW, workflowId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// getWorkflowRhythms returns a workflow's rhythms with their link data, in position order
func getWorkflowRhythms(ctx context.Context, tx *sql.Tx, workflowId uuid.UUID) ([]model.WorkflowRhythm, error) {
	rows, err := tx.QueryContext(ctx, GET_WORKFLOW_RHYTHMS, workflowId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rhythms := make([]model.WorkflowRhythm, 0)
	for rows.Next() {
		var rhythm model.WorkflowRhythm
		err := rows.Scan(
			&rhythm.Id,
			&rhythm.Bpm,
			&rhythm.Beats,
			&rhythm.Subdivision,
			pq.Array(&rhythm.State),
			&rhythm.IsPoly,
			&rhythm.PolyBeats,
			&rhythm.PolySubdivision,
			pq.Array(&rhythm.PolyState),
			&rhythm.UserId,
			&rhythm.Level,
			&rhythm.Name,
			&rhythm.Description,
			&rhythm.Sounds,
			&rhythm.PolySounds,
			&rhythm.CreatedAt,
			&rhythm.UpdatedAt,
			&rhythm.Subdivisions,
			&rhythm.PolySubdivisions,
			&rhythm.Section,
			&rhythm.Measures,
			&rhythm.Position,
		)
		if err != nil {
			return nil, err
		}

		if rhythm.IsPoly {
			rhythm.RhythmType = model.Polyrhythm
		}

		rhythms = append(rhythms, rhythm)
	}

	return rhythms, rows.Err()
}

func getWorkflowSections(ctx context.Context, tx *sql.Tx, workflowId uuid.UUID) ([]model.WorkflowSection, error) {
	rows, err := tx.QueryContext(ctx, GET_WORKFLOW_SECTIONS, workflowId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections := make([]model.WorkflowSection, 0)
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

		sections = append(sections, section)
	}

	return sections, rows.Err()
}

// toWorkflowRhythms pairs each saved rhythm with its block's section and position
func toWorkflowRhythms(blocks []model.CreateWorkflowRhythmInput, rhythms []model.Rhythm) []model.WorkflowRhythm {
	workflowRhythms := make([]model.WorkflowRhythm, len(rhythms))
	for i, rhythm := range rhythms {
		workflowRhythms[i] = model.WorkflowRhythm{
			Rhythm:   rhythm,
			Section:  blocks[i].Section,
			Position: i,
		}
	}

	return workflowRhythms
}

// validateSections checks section ids are unique and that every block
// section references one of them
func validateSections(sections []model.CreateSectionPayload, blockSections []*uuid.UUID) error {
	sectionIds := make(map[uuid.UUID]struct{}, len(sections))
	for _, s := range sections {
		if _, ok := sectionIds[s.SectionId]; ok {
			return ErrDuplicateSection
		}
		sectionIds[s.SectionId] = struct{}{}
	}

	for _, section := range blockSections {
		if section == nil {
			continue
		}
		if _, ok := sectionIds[*section]; !ok {
			return ErrUnknownSection
		}
	}

	return nil
}

func getWorkflowRhythmIds(ctx context.Context, tx *sql.Tx, workflowId uuid.UUID) (map[uuid.UUID]struct{}, error) {
	rows, err := tx.QueryContext(ctx, GET_WORKFLOW_RHYTHM_IDS, workflowId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make(map[uuid.UUID]struct{})
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = struct{}{}
	}

	return ids, rows.Err()
}

// saveSections runs CREATE_SECTIONS or UPSERT_SECTIONS and returns the saved rows
func saveSections(ctx context.Context, tx *sql.Tx, query string, workflowId uuid.UUID, sections []model.CreateSectionPayload) ([]model.WorkflowSection, error) {
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
		query,
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

// saveWorkflowRhythms runs CREATE_WORKFLOW_RHYTHMS or UPSERT_WORKFLOW_RHYTHMS,
// using the order of blocks as each link's position
func saveWorkflowRhythms(ctx context.Context, tx *sql.Tx, query string, workflowId uuid.UUID, blocks []model.CreateWorkflowRhythmInput, rhythms []model.Rhythm) error {
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
		query,
		workflowId,
		pq.Array(rhythmIds),
		pq.Array(sectionIds),
		pq.Array(measures),
		pq.Array(positions),
	)
	return err
}
