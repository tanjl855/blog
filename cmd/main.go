package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/index", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Success",
			"data": struct {
				Data string `json:"data"`
			}{
				Data: "here is index msg",
			},
		})
	})
	r.Run()
}
