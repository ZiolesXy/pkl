package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, errs interface{}) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
		Errors:  errs,
	})
}