package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

// GetInitLogger 简单初始化日志系统，后续优化可带上filePath
func GetInitLogger(infoFileName, warnFileName string) *zap.SugaredLogger {
	encoder := getEncoder()
	// 两个判断日志等级的interface
	// warn level以下属于info
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})
	// warn level及以上属于warn
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	infoWriter := getLogWriter(infoFileName)
	warnWriter := getLogWriter(warnFileName)
	// 创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
	)
	loggers := zap.New(core, zap.AddCaller())
	return loggers.Sugar()
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
