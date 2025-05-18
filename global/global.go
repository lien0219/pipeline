package global

import (
	"gin_pipeline/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	Logger *zap.SugaredLogger
	DB     *gorm.DB
	Redis  *redis.Client
)
