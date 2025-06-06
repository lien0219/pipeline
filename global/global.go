package global

import (
	"gin_pipeline/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Config config.Configuration
	Log    *zap.SugaredLogger
	VP     *viper.Viper
)
