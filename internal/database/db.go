package database

import (
	"fmt"
	"log"
	"time"

	"nalakarsa/internal/config"
	"nalakarsa/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSslMode,
		cfg.DBTimeZone,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.Env == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pooling properties
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection successfully established.")

	// Enable UUID extension in Postgres (required for uuid_generate_v4)
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}

	// Auto Migration
	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		// Core
		&model.User{},
		&model.Profile{},
		&model.RefreshToken{},

		// Discussions
		&model.Discussion{},
		&model.DiscussionReply{},
		&model.DiscussionVote{},

		// Connections
		&model.Connection{},

		// Projects
		&model.Project{},
		&model.ProjectMember{},
		&model.ProjectApplication{},
		&model.ProjectMilestone{},

		// Chat
		&model.Conversation{},
		&model.ConversationMember{},
		&model.Message{},

		// Notifications
		&model.Notification{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration completed successfully.")

	return db, nil
}
