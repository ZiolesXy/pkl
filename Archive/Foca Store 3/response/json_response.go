package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func SuccessList(c *gin.Context, message string, entries interface{}) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": message,
		"data": gin.H{
			"entries": entries,
		},
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
		"data":    nil,
	})
}