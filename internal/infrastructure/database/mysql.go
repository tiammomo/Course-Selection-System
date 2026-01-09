package database

import (
	"fmt"

	"course_select/internal/config"
	"course_select/internal/domain/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// Init 初始化数据库
func Init(cfg *config.DatabaseConfig) error {
	var err error
	db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 自动迁移
	if err := migrate(db); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}

// migrate 执行数据库迁移
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Member{},
		&model.Course{},
		&model.Bind{},
		&model.Choice{},
	)
}

// Get 获取数据库实例
func Get() *gorm.DB {
	return db
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
