package database

import (
	"time"

	"trade_company/internal/config"
	"trade_company/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config, _ any) (*gorm.DB, error) {
	logMode := logger.Info
	if cfg.AppEnv == "production" {
		logMode = logger.Warn
	}
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN()), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Listing{},
		&models.Image{},
		&models.Favorite{},
		&models.Message{},
		&models.Transaction{},
	)
}
