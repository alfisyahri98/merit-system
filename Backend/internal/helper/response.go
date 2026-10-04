package helper

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"data":    data,
	})
}

func Fail(c *gin.Context, code int, msg string, details interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"success": false,
		"error":   msg,
		"details": details,
	})
}
