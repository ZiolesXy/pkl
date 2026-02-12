package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, JSONResponse{Status: "success", Message: message, Data: data})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, JSONResponse{Status: "success", Message: message, Data: data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, JSONResponse{Status: "error", Message: message, Data: gin.H{}})
}
