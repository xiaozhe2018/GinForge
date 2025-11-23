package db

import (
	"gorm.io/gorm"

	"goweb/pkg/config"
	"goweb/pkg/logger"
)

// New 创建数据库连接
func New(cfg *config.Config) (*gorm.DB, error) {
	log := logger.New(
		"database",
		cfg.GetString("log.level"),
		cfg.GetString("log.output"),
		cfg.GetString("log.dir"),
	)
	manager := NewManager(cfg, log)

	if err := manager.Connect(); err != nil {
		return nil, err
	}

	return manager.GetDB(), nil
}
