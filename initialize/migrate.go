package initialize

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
)

// 数据库迁移
func Migrate() error {
	// 自动迁移数据库表结构
	err := global.DB.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		global.Logger.Errorf("数据库迁移失败: %v", err)
		return err
	}

	global.Logger.Info("数据库迁移成功")
	return nil
}
