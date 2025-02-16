package handler

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Errors string `json:"errors"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

