package model

import "github.com/gin-gonic/gin"

// custom error response
type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func NewLinksError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
