package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Logger *zap.SugaredLogger
)

// SetupLogger 创建logger
func SetupLogger() {
	infoFileName := viper.GetString("info_file_name")
	warnFileName := viper.GetString("warn_file_name")

	Logger = GetInitLogger(infoFileName, warnFileName)
	defer Logger.Sync()
}
