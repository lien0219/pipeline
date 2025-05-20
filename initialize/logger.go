package initialize

import (
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	_ "path"
	"time"
)

// InitLogger 初始化日志
func InitLogger() {
	if ok := utils.CreateDir(global.Config.Log.Director); !ok {
		fmt.Printf("创建日志目录 %v 失败！\n", global.Config.Log.Director)
	}

	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/server_debug.log", global.Config.Log.Director), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_info.log", global.Config.Log.Director), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_warn.log", global.Config.Log.Director), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_error.log", global.Config.Log.Director), errorPriority),
	}

	logger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	if global.Config.Log.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	// 赋值给全局变量
	global.Log = logger.Sugar()
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// customTimeEncoder 自定义日志输出时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(global.Config.Log.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if global.Config.Log.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) zapcore.Core {
	writer := getWriteSyncer(fileName)
	return zapcore.NewCore(getEncoder(), writer, level)
}

// getWriteSyncer 获取zapcore.WriteSyncer
func getWriteSyncer(fileName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    global.Config.Log.MaxSize,
		MaxBackups: global.Config.Log.MaxBackups,
		MaxAge:     global.Config.Log.MaxAge,
		Compress:   true,
	}

	if global.Config.Log.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
