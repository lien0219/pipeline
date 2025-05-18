package initialize

import (
	"context"
	"fmt"
	"gin_pipeline/config"
	"gin_pipeline/global"
	"github.com/go-redis/redis/v8"
)

// Redis 初始化Redis连接
func Redis() error {
	redisConfig := config.Get().Redis

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}

	global.Redis = client
	return nil
}
