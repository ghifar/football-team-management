package test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func MakeRequest(router *gin.Engine, method string, URI string, body io.Reader) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, URI, body)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	return response
}

func Router(URI string, handler func(c *gin.Context), method string) *gin.Engine {
	router := gin.New()
	router.Handle(method, URI, handler)
	return router
}
