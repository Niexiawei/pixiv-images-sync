package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"pixivImages/config"
	"pixivImages/utils"
)

var (
	Logger       *zap.SugaredLogger
	loggerConfig config.Logger
)

func InitLogger() {
	loggerConfig = config.Get().Logger
	if ok, _ := utils.PathExists(loggerConfig.Path); !ok {
		_ = os.MkdirAll(loggerConfig.Path, 0777)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(stdoutEncoder(), stdoutWriter(), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder(), fileWriter(), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zap.InfoLevel
		})),
	)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func stdoutEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func fileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func fileWriter() zapcore.WriteSyncer {
	filePath, _ := filepath.Abs(loggerConfig.Path + "/logger.log")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,                // 文件位置
		MaxSize:    loggerConfig.MaxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     loggerConfig.MaxAge,     // 保留旧文件的最大天数
		MaxBackups: loggerConfig.MaxBackups, // 保留旧文件的最大个数
		Compress:   false,                   // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

func stdoutWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}
