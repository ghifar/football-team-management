package handlers

import (
	"context"
	"football-team-management/internal/domain"
	"football-team-management/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	repo usecases.PlayerRepository
}

func NewPlayerHandler(repo usecases.PlayerRepository) *PlayerHandler {
	return &PlayerHandler{repo: repo}
}

func (h *PlayerHandler) Register(c *gin.Context) {
	var player domain.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Register(context.Background(), player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) Update(c *gin.Context) {
	name := c.Param("playerName")
	var player domain.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Update(context.Background(), name, player); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Delete(c *gin.Context) {
	name := c.Param("playerName")
	if err := h.repo.Delete(context.Background(), name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "player deleted"})
}

func (h *PlayerHandler) List(c *gin.Context) {
	players, err := h.repo.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) ListByTeam(c *gin.Context) {
	teamName := c.Param("teamName")
	players, err := h.repo.ListByTeam(context.Background(), teamName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) GetByName(c *gin.Context) {
	name := c.Param("playerName")
	player, err := h.repo.GetByName(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Restore(c *gin.Context) {
	name := c.Param("playerName")
	if err := h.repo.Restore(context.Background(), name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "player restored"})
}
