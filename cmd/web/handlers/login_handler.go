package handlers

import (
	"football-team-management/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler interface {
	Handle(c *gin.Context)
}

type loginHandlerImpl struct {
	authService usecases.AuthService
}

func NewLoginHandlerImpl(authService usecases.AuthService) *loginHandlerImpl {
	return &loginHandlerImpl{authService: authService}
}

func (l *loginHandlerImpl) Handle(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := l.authService.GenerateToken(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}
