package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rhythmapi/handler/rhythm"
	"rhythmapi/model"
	"rhythmapi/request"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func (r *RhythmHandler) GetRhythmById(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	foundRhythm, err := r.repo.GetRhythmById(c, parsedId)
	if err != nil {
		log := fmt.Sprintf("Error fetching rhythm: %s", err.Error())
		respondErr(c, log, rhythm.ErrRhythmById)
		return
	}

	log := fmt.Sprintf("found rhythm %v for user %v", foundRhythm.Id.String(), foundRhythm.UserId.String())
	respondSuccessContent(c, log, gin.H{"rhythm": foundRhythm})
}

func (r *RhythmHandler) DeleteRhythmById(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	deleted, err := r.repo.DeleteRhythmById(c, parsedId)
	if err != nil || !deleted {
		if errors.Is(err, sql.ErrNoRows) {
			log := fmt.Sprintf("no rhythm found by id: %v", parsedId.String())
			respondErr(c, log, rhythm.ErrDeleteRhythmById)
			return
		}

		log := fmt.Sprintf("failed to delete rhythm by id: %v", parsedId.String())
		respondErr(c, log, rhythm.ErrDeleteRhythmById)
		return
	}

	log := fmt.Sprintf("successfully deleted rhythm %v", parsedId.String())
	responseSuccessNoContent(c, log)
}

func (r *RhythmHandler) UpdateRhythmById(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var rhythmToUpdate model.RhythmInputPoly
	if err = ctx.ShouldBindJSON(&rhythmToUpdate); err != nil {
		log := fmt.Sprintf("failed to parse rhythm payload: %v", err.Error())
		respondErr(ctx, log, rhythm.ErrUpdateRhythmById)
		return
	}

	isValid := request.ValidatePolyRhythm(rhythmToUpdate)
	if !isValid {
		log := "poly rhythms must contain associated poly fields"
		respondErr(ctx, log, rhythm.ErrUpdateRhythmById)
		return
	}

	updatedRhythm, err := r.repo.UpdateRhythmById(ctx, parsedId, rhythmToUpdate)
	if err != nil {
		log := fmt.Sprintf("failed to update rhythm %v: %v", parsedId.String(), err.Error())
		respondErr(ctx, log, rhythm.ErrUpdateRhythmById)
		return
	}

	log := fmt.Sprintf("successfully update rhythm %v", parsedId.String())
	respondSuccessContent(ctx, log, gin.H{"rhythm": updatedRhythm})
}

func (r *RhythmHandler) GetRhythms(ctx *gin.Context) {
	var rhythmRequest model.GetRhythmsRequest

	offset := ctx.Query("offset")
	if offset == "" {
		rhythmRequest.Offset = 1
	} else {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			log := fmt.Sprintf("invalid offset: %v", err.Error())
			respondErr(ctx, log, rhythm.ErrGetRhythms)
			return
		}

		rhythmRequest.Offset = parsedOffset
	}

	rhythms, total, err := r.repo.GetRhythms(ctx, rhythmRequest)
	if err != nil {
		log := fmt.Sprintf("failed to get rhythms at offset %d: %v", rhythmRequest.Offset, err.Error())
		respondErr(ctx, log, rhythm.ErrGetRhythms)
		return
	}

	respondSuccessContent(ctx, "successfully got rhythms", gin.H{"rhythms": rhythms, "total": total})
}

func (r *RhythmHandler) GetSubdivisionTypes(ctx *gin.Context) {
	subdivisionTypes, err := r.repo.GetSubdivisionTypes(ctx)
	if err != nil {
		log.Println("Error fetching rhythm:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "rhythm not found"})
		return
	}

	ctx.JSON(http.StatusOK, subdivisionTypes)
}

func (r *RhythmHandler) GetRhythmLevels(ctx *gin.Context) {
	rhythmLevels, err := r.repo.GetRhythmLevels(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "rhythm levels not found"})
		return
	}

	ctx.JSON(http.StatusOK, rhythmLevels)
}

// func (r *RhythmHandler) CreateTags(ctx *gin.Context) {
// 	var tags model.CreateTagInput

// 	if err := ctx.ShouldBindBodyWithJSON(&tags); err != nil {
// 		respondErr(ctx, "failed to parse tag body", rhythm.ErrTagCreation)
// 		return
// 	}

// }

func (r *RhythmHandler) CreateRhtyhm(ctx *gin.Context) {
	var discriminator model.RhythmDiscriminator

	err := ctx.ShouldBindBodyWith(&discriminator, binding.JSON)
	if err != nil {
		respondErr(ctx, "failed to parse rhythm type", rhythm.ErrRhythmCreateError)
		return
	}

	switch discriminator.RhythmType {
	case model.Monorhythm:
		var input model.RhythmInputMono

		if err := ctx.ShouldBindBodyWith(&input, binding.JSON); err != nil {
			respondErr(ctx, "failed to parse mono rhythm", rhythm.ErrRhythmCreateError)
			return
		}

		r, err := r.repo.CreateMonoRhythm(ctx, input)
		if err != nil {
			respondErr(ctx, err.Error(), rhythm.ErrRhythmCreateError)
			return
		}

		log := fmt.Sprintf("created mono rhythm for: %d - rhythm: %d", r.UserId, r.Id)
		respondSuccessContent(ctx, log, gin.H{"rhythm": r})

	case model.Polyrhythm:
		var input model.RhythmInputPoly

		if err := ctx.ShouldBindBodyWith(&input, binding.JSON); err != nil {
			respondErr(ctx, "failed to parse poly rhythm"+err.Error(), rhythm.ErrRhythmCreateError)
			return
		}

		r, err := r.repo.CreatePolyRhythm(ctx, input)
		if err != nil {
			respondErr(ctx, err.Error(), rhythm.ErrRhythmCreateError)
			return
		}

		log := fmt.Sprintf("created poly rhythm for: %d - rhythm: %d", r.UserId, r.Id)
		respondSuccessContent(ctx, log, gin.H{"rhythm": r})

	default:
		respondErr(ctx, "failed to parse rhythm - invalid rhythm type", rhythm.ErrRhythmCreateError)
		return
	}
}
