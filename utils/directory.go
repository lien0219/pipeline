package utils

import (
	"os"
	"path/filepath"
)

// CreateDir 创建目录
func CreateDir(path string) bool {
	// 检查目录是否存在
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	// 创建目录
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

// PathExists 判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetCurrentDirectory 获取当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}
