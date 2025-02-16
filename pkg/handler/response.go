package handler

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Errors string `json:"errors"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type userResponse struct {
	ID       string `json:"id" example:"+79161234567"`
	Username string `json:"username" example:"some.mail@gmail.com"`
	Coins    int    `json:"coins" example:"25"`
}

type createResponse struct {
	ID uuid `json:"id"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
