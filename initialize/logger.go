package initialize

import (
	"gin_pipeline/config"
	"gin_pipeline/global"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义颜色常量
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// 自定义日志级别颜色编码器
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var levelStr string

	switch level {
	case zapcore.DebugLevel:
		levelStr = colorBlue + "DEBUG" + colorReset
	case zapcore.InfoLevel:
		levelStr = colorCyan + "INFO" + colorReset
	case zapcore.WarnLevel:
		levelStr = colorGreen + "WARN" + colorReset // 绿色警告
	case zapcore.ErrorLevel:
		levelStr = colorRed + "ERROR" + colorReset // 红色错误
	case zapcore.DPanicLevel:
		levelStr = colorPurple + "DPANIC" + colorReset
	case zapcore.PanicLevel:
		levelStr = colorPurple + "PANIC" + colorReset
	case zapcore.FatalLevel:
		levelStr = colorRed + "FATAL" + colorReset
	default:
		levelStr = colorWhite + level.String() + colorReset
	}

	enc.AppendString(levelStr)
}

// Logger 初始化日志
func Logger() error {
	logConfig := config.Get().Log

	// 确保日志目录存在
	dir := logConfig.Filename[:len(logConfig.Filename)-len("/app.log")]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 日志分割
	hook := lumberjack.Logger{
		Filename:   logConfig.Filename,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
	}

	// 设置日志级别
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 文件输出编码器 - 使用JSON格式
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 文件中使用小写字母级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// 控制台输出编码器 - 使用彩色输出
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder, // 使用自定义彩色级别编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// 创建核心 - 同时写入文件和控制台
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(&hook), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	// 创建Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	global.Logger = logger.Sugar()

	return nil
}
