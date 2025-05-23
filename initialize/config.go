package initialize

import (
	"fmt"
	"gin_pipeline/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig() {
	var config string
	// 判断命令行参数是否有配置文件
	if len(os.Args) > 1 {
		config = os.Args[1]
	} else {
		config = "config.yaml"
		fmt.Printf("未指定配置文件，使用默认配置文件: %s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s", err))
	}

	// 监控配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生变化...")
		if err := v.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 将配置赋值给全局变量
	if err := v.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}

	// 赋值给全局viper
	global.VP = v
}
