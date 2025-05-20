package initialize

import (
	"context"
	"gin_pipeline/global"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// InitRedis 初始化Redis
func InitRedis() {
	redisConfig := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Log.Error("Redis连接失败", zap.Any("err", err))
		return
	}
	global.Redis = client
}
