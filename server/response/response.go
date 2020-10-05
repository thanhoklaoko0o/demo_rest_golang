package response

import (
	"github.com/gin-gonic/gin"
)

// JSON displays a json response message with data
func JSON(context *gin.Context, statusCode int, data interface{}) {
	context.JSON(statusCode, gin.H{
		"data": data,
	})
}

// ERROR displays a json response message with an error
func ERROR(context *gin.Context, statusCode int) {
	context.JSON(statusCode, gin.H{
		"error": statusCode,
	})
}
