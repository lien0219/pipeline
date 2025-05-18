package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var (
	config *Config
	once   sync.Once
)

// 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// 服务器配置
type ServerConfig struct {
	Mode         string `yaml:"mode"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

// 数据库配置
type DatabaseConfig struct {
	Driver      string `yaml:"driver"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	Charset     string `yaml:"charset"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxOpen     int    `yaml:"maxOpen"`
	MaxLifetime int    `yaml:"maxLifetime"`
	ShowSQL     bool   `yaml:"showSql"`
}

// Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

// 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

// JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expireTime"` // 过期时间（小时）
}

// 初始化配置
func Init(configPath string) error {
	var err error
	once.Do(func() {
		var configData []byte
		configData, err = ioutil.ReadFile(configPath)
		if err != nil {
			return
		}

		config = &Config{}
		err = yaml.Unmarshal(configData, config)
	})

	return err
}

// Get 获取配置
func Get() *Config {
	if config == nil {
		panic("Configuration not initialized")
	}
	return config
}
