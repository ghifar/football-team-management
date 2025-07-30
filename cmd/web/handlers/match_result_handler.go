package handlers

import (
	"context"
	"football-team-management/internal/domain"
	"football-team-management/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchResultHandler struct {
	repo usecases.MatchResultRepository
}

func NewMatchResultHandler(repo usecases.MatchResultRepository) *MatchResultHandler {
	return &MatchResultHandler{repo: repo}
}

func (h *MatchResultHandler) Register(c *gin.Context) {
	var resultReq domain.MatchResultRequest
	if err := c.ShouldBindJSON(&resultReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := resultReq.ToMatchResult()
	if err := h.repo.Register(context.Background(), *result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := result.ToMatchResultResponse()
	c.JSON(http.StatusCreated, response)
}

func (h *MatchResultHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match result id"})
		return
	}

	var resultReq domain.MatchResultRequest
	if err := c.ShouldBindJSON(&resultReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := resultReq.ToMatchResult()
	if err := h.repo.Update(context.Background(), id, *result); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := result.ToMatchResultResponse()
	response.ID = id
	c.JSON(http.StatusOK, response)
}

func (h *MatchResultHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match result id"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match result deleted"})
}

func (h *MatchResultHandler) List(c *gin.Context) {
	results, err := h.repo.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*domain.MatchResultResponse
	for _, result := range results {
		responses = append(responses, result.ToMatchResultResponse())
	}
	c.JSON(http.StatusOK, responses)
}

func (h *MatchResultHandler) GetByMatchID(c *gin.Context) {
	matchIDStr := c.Param("matchID")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	result, err := h.repo.GetByMatchID(context.Background(), matchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := result.ToMatchResultResponse()
	c.JSON(http.StatusOK, response)
}

func (h *MatchResultHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match result id"})
		return
	}

	result, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := result.ToMatchResultResponse()
	c.JSON(http.StatusOK, response)
}

func (h *MatchResultHandler) Restore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match result id"})
		return
	}

	if err := h.repo.Restore(context.Background(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match result restored"})
}
