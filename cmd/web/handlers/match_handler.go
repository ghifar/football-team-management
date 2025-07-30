package handlers

import (
	"context"
	"football-team-management/internal/domain"
	"football-team-management/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	repo usecases.MatchRepository
}

func NewMatchHandler(repo usecases.MatchRepository) *MatchHandler {
	return &MatchHandler{repo: repo}
}

func (h *MatchHandler) Register(c *gin.Context) {
	var matchReq domain.MatchRequest
	if err := c.ShouldBindJSON(&matchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := matchReq.ToMatch()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. Use YYYY-MM-DD"})
		return
	}

	if err := h.repo.Register(context.Background(), *match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := match.ToMatchResponse()
	c.JSON(http.StatusCreated, response)
}

func (h *MatchHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	var matchReq domain.MatchRequest
	if err := c.ShouldBindJSON(&matchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := matchReq.ToMatch()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. Use YYYY-MM-DD"})
		return
	}

	if err := h.repo.Update(context.Background(), id, *match); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := match.ToMatchResponse()
	response.ID = id
	c.JSON(http.StatusOK, response)
}

func (h *MatchHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match deleted"})
}

func (h *MatchHandler) List(c *gin.Context) {
	matches, err := h.repo.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*domain.MatchResponse
	for _, match := range matches {
		responses = append(responses, match.ToMatchResponse())
	}
	c.JSON(http.StatusOK, responses)
}

func (h *MatchHandler) ListByTeam(c *gin.Context) {
	teamName := c.Param("teamName")
	matches, err := h.repo.ListByTeam(context.Background(), teamName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*domain.MatchResponse
	for _, match := range matches {
		responses = append(responses, match.ToMatchResponse())
	}
	c.JSON(http.StatusOK, responses)
}

func (h *MatchHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	match, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := match.ToMatchResponse()
	c.JSON(http.StatusOK, response)
}

func (h *MatchHandler) Restore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	if err := h.repo.Restore(context.Background(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match restored"})
}
