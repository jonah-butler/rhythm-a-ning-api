package repository

import (
	"database/sql"
	"rhythmapi/model"
	"rhythmapi/request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RhythmRepository struct {
	db *sql.DB
}

func NewRhythmRepository(db *sql.DB) IRhythmRepository {
	return &RhythmRepository{db: db}
}

func (r *RhythmRepository) GetRhythmById(ctx *gin.Context, id uuid.UUID) (model.Rhythm, error) {
	var rhythm model.Rhythm

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return rhythm, err
	}

	err = r.db.QueryRowContext(
		ctx,
		GET_RHYTHM_BY_ID,
		userId,
		id,
	).Scan(
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
		&rhythm.CreatedAt,
		&rhythm.UpdatedAt,
		&rhythm.Sounds,
		&rhythm.PolySounds,
	)
	if err != nil {
		return rhythm, err
	}

	if rhythm.IsPoly {
		rhythm.RhythmType = model.Polyrhythm
	}

	return rhythm, nil
}

func (r *RhythmRepository) DeleteRhythmById(ctx *gin.Context, rhythmId uuid.UUID) (bool, error) {
	var deleted bool

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return deleted, err
	}

	err = r.db.QueryRowContext(
		ctx,
		DELETE_RHYTHM_BY_ID,
		rhythmId,
		userId,
	).Scan(&deleted)
	if err != nil {
		return deleted, err
	}

	if !deleted {
		return deleted, sql.ErrNoRows
	}

	return deleted, nil
}

func (r *RhythmRepository) GetRhythms(ctx *gin.Context, rhythmsRequest model.GetRhythmsRequest) ([]model.Rhythm, int, error) {
	var rhythms []model.Rhythm

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return rhythms, 0, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		GET_RHYTHMS,
		userId,
		rhythmsRequest.Offset,
	)
	if err != nil {
		return rhythms, 0, err
	}

	defer rows.Close()

	var total int
	for rows.Next() {
		var rhythm model.Rhythm

		err = rows.Scan(
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
			&total,
		)
		if err != nil {
			return rhythms, total, err
		}

		rhythms = append(rhythms, rhythm)
	}

	return rhythms, total, nil
}

func (r *RhythmRepository) UpdateRhythmById(ctx *gin.Context, rhythmId uuid.UUID, update model.RhythmInputPoly) (model.Rhythm, error) {
	var rhythm model.Rhythm

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return rhythm, err
	}

	var isPoly bool
	if update.RhythmType == model.Polyrhythm {
		isPoly = true
	}

	if update.RhythmType == model.Monorhythm {
		update.PolyState = nil
		update.PolyBeats = nil
		update.PolySubdivision = nil
	}

	err = r.db.QueryRowContext(
		ctx,
		UPDATE_RHYTHM_BY_ID,
		update.Bpm,
		update.Beats,
		update.Subdivision,
		pq.Array(update.State),
		isPoly,
		update.PolyBeats,
		update.PolySubdivision,
		pq.Array(update.PolyState),
		update.Level,
		update.Name,
		update.Description,
		update.Sounds,
		update.PolySounds,
		rhythmId,
		userId,
	).Scan(
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
		&rhythm.Name,
		&rhythm.Description,
		&rhythm.Level,
		&rhythm.Sounds,
		&rhythm.PolySounds,
		&rhythm.CreatedAt,
		&rhythm.UpdatedAt,
	)
	if err != nil {
		return rhythm, err
	}

	if rhythm.IsPoly {
		rhythm.RhythmType = model.Polyrhythm
	}

	return rhythm, nil
}

func (r *RhythmRepository) GetSubdivisionTypes(ctx *gin.Context) ([]model.SubdivisionType, error) {
	var subdivisionTypes []model.SubdivisionType

	rows, err := r.db.QueryContext(ctx, GET_SUBDIVISION_TYPES)
	if err != nil {
		return subdivisionTypes, err
	}

	defer rows.Close()

	for rows.Next() {
		var subdivisionType model.SubdivisionType

		err = rows.Scan(&subdivisionType.SubdivisionId, &subdivisionType.Name)
		if err != nil {
			return subdivisionTypes, err
		}

		subdivisionTypes = append(subdivisionTypes, subdivisionType)
	}

	return subdivisionTypes, nil
}

func (r *RhythmRepository) GetRhythmLevels(ctx *gin.Context) ([]model.RhythmLevel, error) {
	var rhythmLevels []model.RhythmLevel

	rows, err := r.db.QueryContext(ctx, GET_RHYTHM_LEVELS)
	if err != nil {
		return rhythmLevels, err
	}

	defer rows.Close()

	for rows.Next() {
		var rhythmLevel model.RhythmLevel

		err = rows.Scan(&rhythmLevel.LevelId, &rhythmLevel.Name)
		if err != nil {
			return rhythmLevels, err
		}

		rhythmLevels = append(rhythmLevels, rhythmLevel)
	}

	return rhythmLevels, nil
}

func (r *RhythmRepository) CreateMonoRhythm(ctx *gin.Context, rhythm model.RhythmInputMono) (model.Rhythm, error) {
	var inserted model.Rhythm

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return inserted, err
	}

	err = r.db.QueryRowContext(
		ctx,
		CREATE_RHYTHM,
		rhythm.Bpm,
		rhythm.Beats,
		rhythm.Subdivision,
		pq.Array(rhythm.State),
		false,
		nil,
		nil,
		nil,
		nil,
		userId,
		rhythm.Name,
		rhythm.Description,
		rhythm.Level,
		rhythm.Sounds,
	).Scan(
		&inserted.Id,
		&inserted.Bpm,
		&inserted.Beats,
		&inserted.Subdivision,
		pq.Array(&inserted.State),
		&inserted.IsPoly,
		&inserted.PolyBeats,
		&inserted.PolySubdivision,
		pq.Array(&inserted.PolyState),
		&inserted.PolySounds,
		&inserted.UserId,
		&inserted.CreatedAt,
		&inserted.UpdatedAt,
		&inserted.Name,
		&inserted.Description,
		&inserted.Level,
		&inserted.Sounds,
	)
	if err != nil {
		return inserted, err
	}

	inserted.RhythmType = model.Monorhythm

	return inserted, nil
}

func (r *RhythmRepository) CreatePolyRhythm(ctx *gin.Context, rhythm model.RhythmInputPoly) (model.Rhythm, error) {
	var inserted model.Rhythm

	userId, err := request.GetUserFromContext(ctx)
	if err != nil {
		return inserted, err
	}

	err = r.db.QueryRowContext(
		ctx,
		CREATE_RHYTHM,
		rhythm.Bpm,
		rhythm.Beats,
		rhythm.Subdivision,
		pq.Array(rhythm.State),
		true,
		rhythm.PolyBeats,
		rhythm.PolySubdivision,
		rhythm.PolyState,
		rhythm.PolySounds,
		userId,
		rhythm.Name,
		rhythm.Description,
		rhythm.Level,
		rhythm.Sounds,
	).Scan(
		&inserted.Id,
		&inserted.Bpm,
		&inserted.Beats,
		&inserted.Subdivision,
		pq.Array(&inserted.State),
		&inserted.IsPoly,
		&inserted.PolyBeats,
		&inserted.PolySubdivision,
		pq.Array(&inserted.PolyState),
		&inserted.PolySounds,
		&inserted.UserId,
		&inserted.CreatedAt,
		&inserted.UpdatedAt,
		&inserted.Name,
		&inserted.Description,
		&inserted.Level,
		&inserted.Sounds,
	)
	if err != nil {
		return inserted, err
	}

	inserted.RhythmType = model.Polyrhythm

	return inserted, nil
}
