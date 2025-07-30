package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler interface {
	Handle(c *gin.Context)
}

type pingHandlerImpl struct {
}

func NewPingHandlerImpl() *pingHandlerImpl {
	return &pingHandlerImpl{}
}

func (p pingHandlerImpl) Handle(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

var _ PingHandler = (*pingHandlerImpl)(nil)
