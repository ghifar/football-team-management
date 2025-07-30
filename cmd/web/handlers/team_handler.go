package handlers

import (
	"context"
	"football-team-management/internal/domain"
	"football-team-management/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	repo usecases.TeamRepository
}

func NewTeamHandler(repo usecases.TeamRepository) *TeamHandler {
	return &TeamHandler{repo: repo}
}

func (h *TeamHandler) Register(c *gin.Context) {
	var team domain.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Register(context.Background(), team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var team domain.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Update(context.Background(), name, team); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.repo.Delete(context.Background(), name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "team deleted"})
}

func (h *TeamHandler) List(c *gin.Context) {
	teams, err := h.repo.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) Restore(c *gin.Context) {
	name := c.Param("name")
	if err := h.repo.Restore(context.Background(), name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "team restored"})
}
