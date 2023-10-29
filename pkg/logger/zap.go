package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

var (
	ErrorLogger *zap.Logger
	InfoLogger  *zap.Logger
)

func InitZapLogger() {
	logPath := viper.GetString("log_path")
	hostName := viper.GetString("host_name")
	ErrorLogger = getNewZap(fmt.Sprintf("%s/%s/%s_%s", logPath, "error", hostName, "error"))
	InfoLogger = getNewZap(fmt.Sprintf("%s/%s/%s_%s", logPath, "info", hostName, "info"))
}

func getNewZap(fileName string) *zap.Logger {
	writeSyncer := getLogWriter(fileName)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	return zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	ioWriter := getWriter(fileName)
	return zapcore.AddSync(ioWriter)
}

// getWriter
// 日志文件按天切割
func getWriter(fileName string) io.Writer {
	// 保存30天内的日志, 每24个小时（整点）分割一次日志
	hook, err := rotatelogs.New(
		fileName+"_%Y%m%d.log",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
