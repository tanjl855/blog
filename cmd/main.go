package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tanjl855/blog/internal/blog/router"
	"github.com/tanjl855/blog/pkg/conf"
	"github.com/tanjl855/blog/pkg/database"
	"github.com/tanjl855/blog/pkg/logger"
	"github.com/tanjl855/blog/pkg/utils"
	"io"
	"net/http"
	"time"
)

func init() {
	utils.InitFlag()
	conf.InitConfig() // 初始化viper读config
	setEnv()
	logger.SetupLogger() // 初始化logger

	database.InitMysql()
}

func setEnv() {
	env := viper.GetString("env")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func shutdownBlog() {
	database.BlogSqlDB.Close()
}

func main() {
	defer shutdownBlog()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", viper.GetString("server_port")),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router.SetupRouter(),
	}

	srv.ListenAndServe()
}
