package http

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestShowSuccessWithData(t *testing.T) {
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		ShowSuccessWithData(c, "test here!")
	})
	r.Run()
}
