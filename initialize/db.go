package initialize

import (
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitDB 初始化数据库
func InitDB() {
	m := global.Config.Mysql
	if m.DbName == "" {
		global.Log.Error("未配置数据库名")
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", m.Username, m.Password, m.Path, m.Port, m.DbName, m.Config)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}

	// 设置日志级别
	var logLevel logger.LogLevel
	switch m.LogMode {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// 创建自定义日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})

	if err != nil {
		global.Log.Error("数据库连接失败", zap.Any("err", err))
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	global.DB = db

	// 自动迁移
	err = db.AutoMigrate(
		&model.User{},
		&model.Pipeline{},
		&model.Stage{},
		&model.Job{},
		&model.PipelineRun{},
		&model.Artifact{},
		&model.Environment{},
		&model.Release{},
		&model.BuildTemplate{},
		&model.DAG{},
		&model.YAMLValidation{},
		&model.YAMLSchema{},
		&model.TemplateCategory{},
		&model.Template{},
		&model.TemplateVersion{},
		&model.Cluster{},
		&model.HPAPolicy{},
		&model.ResourceQuota{},
		&model.TenantResourceRequest{},
		&model.ResourceReport{},
	)
	if err != nil {
		global.Log.Error("自动迁移失败", zap.Any("err", err))
		return
	}
	global.Log.Info("自动迁移成功")
}
