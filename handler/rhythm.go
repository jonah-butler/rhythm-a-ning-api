package handler

import (
	"fmt"
	"log"
	"net/http"
	"rhythmapi/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RhythmHandler struct {
	repo repository.IRhythmRepository
}

func NewRhythmHandler(repo repository.IRhythmRepository) *RhythmHandler {
	return &RhythmHandler{repo: repo}
}

func (r *RhythmHandler) GetRhythmById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rhythm, err := r.repo.FindById(idInt)
	if err != nil {
		log.Println("Error fetching rhythm:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "rhythm not found"})
		return
	}

	c.JSON(http.StatusOK, rhythm)
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
	fmt.Println("getting rhythm levels")
	rhythmLevels, err := r.repo.GetRhythmLevels(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "rhythm levels not found"})
		return
	}

	ctx.JSON(http.StatusOK, rhythmLevels)
}
