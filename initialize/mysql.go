package initialize

import (
	"fmt"
	"gin_pipeline/config"
	"gin_pipeline/global"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// MySQL 初始化MySQL连接
func MySQL() error {
	dbConfig := config.Get().Database

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.Charset,
	)

	// 设置GORM日志级别
	var logLevel logger.LogLevel
	if dbConfig.ShowSQL {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return err
	}

	// 获取底层SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return err
	}

	global.DB = db
	return nil
}
