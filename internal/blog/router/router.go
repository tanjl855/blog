package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tanjl855/blog/internal/blog/controller"
	response "github.com/tanjl855/blog/internal/blog/http"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	blogCors := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	})

	if viper.GetString("env") != "production" {
		pprof.Register(r)
	}

	r.Use()

	r.Use(gin.Recovery(), blogCors)

	r.NoRoute(func(c *gin.Context) {
		rsp := gin.H{"errCode": response.GatewayNotFound, "message": response.ResponseMessages[response.GatewayNotFound]}
		c.JSON(http.StatusOK, rsp)
	})

	r.GET("/blog/v1/configs", controller.Configs)

	return r
}
