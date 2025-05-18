package main

import (
	"fmt"
	"gin_admin_api/config"
)

func main() {
	fmt.Println("测试启动")
	fmt.Println("系统配置：", config.Config.System)
}
