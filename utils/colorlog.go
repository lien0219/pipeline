package utils

import (
	"fmt"
	"gin_pipeline/global"
	"runtime"
	"strings"
)

const (
	greenBg   = "\033[97;42m"
	whiteBg   = "\033[90;47m"
	yellowBg  = "\033[90;43m"
	redBg     = "\033[97;41m"
	blueBg    = "\033[97;44m"
	magentaBg = "\033[97;45m"
	cyanBg    = "\033[97;46m"
	green     = "\033[32m"
	white     = "\033[37m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	reset     = "\033[0m"
)

// ColorfulPrint 彩色打印
func ColorfulPrint(level, format string, args ...interface{}) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	// 提取文件名
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	// 根据日志级别选择颜色
	var levelColor, msgColor string
	switch strings.ToUpper(level) {
	case "INFO":
		levelColor = blueBg
		msgColor = blue
	case "DEBUG":
		levelColor = magentaBg
		msgColor = magenta
	case "WARN":
		levelColor = yellowBg
		msgColor = yellow
	case "ERROR":
		levelColor = redBg
		msgColor = red
	case "SUCCESS":
		levelColor = greenBg
		msgColor = green
	default:
		levelColor = whiteBg
		msgColor = white
	}

	// 格式化消息
	message := fmt.Sprintf(format, args...)

	// 输出彩色日志
	fmt.Printf("%s %s %s [%s:%d] %s%s%s\n",
		levelColor, level, reset,
		file, line,
		msgColor, message, reset)

	// 同时使用zap记录日志
	switch strings.ToUpper(level) {
	case "INFO":
		global.Log.Infof("[%s:%d] %s", file, line, message)
	case "DEBUG":
		global.Log.Debugf("[%s:%d] %s", file, line, message)
	case "WARN":
		global.Log.Warnf("[%s:%d] %s", file, line, message)
	case "ERROR":
		global.Log.Errorf("[%s:%d] %s", file, line, message)
	case "SUCCESS":
		global.Log.Infof("[SUCCESS] [%s:%d] %s", file, line, message)
	default:
		global.Log.Infof("[%s:%d] %s", file, line, message)
	}
}

// Info 输出信息日志
func Info(format string, args ...interface{}) {
	ColorfulPrint("INFO", format, args...)
}

// Debug 输出调试日志
func Debug(format string, args ...interface{}) {
	ColorfulPrint("DEBUG", format, args...)
}

// Warn 输出警告日志
func Warn(format string, args ...interface{}) {
	ColorfulPrint("WARN", format, args...)
}

// Error 输出错误日志
func Error(format string, args ...interface{}) {
	ColorfulPrint("ERROR", format, args...)
}

// Success 输出成功日志
func Success(format string, args ...interface{}) {
	ColorfulPrint("SUCCESS", format, args...)
}
