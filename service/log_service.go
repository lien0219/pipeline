package service

import (
	"bufio"
	"fmt"
	"gin_pipeline/global"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LogService 日志服务
type LogService struct{}

// LogFile 日志文件信息
type LogFile struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

// GetLogFiles 获取日志文件列表
func (s *LogService) GetLogFiles() ([]LogFile, error) {
	logDir := global.Config.Log.Director
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, err
	}

	var logFiles []LogFile
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" {
			info, _ := file.Info()
			modified := info.ModTime().Format("2006-01-02 15:04:05")
			logFiles = append(logFiles, LogFile{
				Name:     file.Name(),
				Size:     info.Size(),
				Modified: modified,
			})
		}
	}

	// 按修改时间排序
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Modified > logFiles[j].Modified
	})

	return logFiles, nil
}

// GetLogContent 读取日志文件内容
func (s *LogService) GetLogContent(fileName string, startLine, limit int) (string, int, error) {
	logDir := global.Config.Log.Director
	filePath := filepath.Join(logDir, fileName)

	// 安全检查（保持原有逻辑）
	absLogDir, err := filepath.Abs(logDir)
	if err != nil {
		return "", 0, err
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", 0, err
	}
	if !strings.HasPrefix(absFilePath, absLogDir) {
		return "", 0, fmt.Errorf("invalid file path")
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	lineCount := 0
	currentLine := 0

	// 读取并收集指定范围的行
	for scanner.Scan() {
		lineCount++
		if currentLine >= startLine && len(lines) < limit {
			lines = append(lines, scanner.Text())
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}

	// 返回内容和总行数
	return strings.Join(lines, "\n"), lineCount, nil
}
