package controller

import "github.com/gin-gonic/gin"

type Response struct {
	Error string `json:"error"`
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Response{Error: message})
}
