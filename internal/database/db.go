package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"nalakarsa/internal/config"
	"nalakarsa/internal/model/connection"
	"nalakarsa/internal/model/conversation"
	"nalakarsa/internal/model/discussion"
	"nalakarsa/internal/model/homepage"
	"nalakarsa/internal/model/notification"
	"nalakarsa/internal/model/project"
	"nalakarsa/internal/model/user"

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

	dbLogger := logger.New(
		log.New(os.Stdout, "", 0),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError:  true,
			Colorful:                  true,
		},
	)

	if cfg.Env == "production" {
		dbLogger = logger.New(
			log.New(os.Stdout, "", 0),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError:  true,
				Colorful:                  false,
			},
		)
	}

	gormConfig := &gorm.Config{
		Logger: dbLogger,
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
	if cfg.DBAutoMigrate {
		log.Println("Running database migrations...")
		err = db.AutoMigrate(
			// Core
			&user.User{},

			&user.RefreshToken{},

			// Discussions
			&discussion.Discussion{},
			&discussion.DiscussionReply{},
			&discussion.DiscussionVote{},

			// Connections
			&connection.Connection{},

			// Projects
			&project.Project{},
			&project.ProjectMember{},
			&project.ProjectApplication{},
			&project.ProjectMilestone{},

			// Chat
			&conversation.Conversation{},
			&conversation.ConversationMember{},
			&conversation.Message{},

			// Notifications
			&notification.Notification{},

			// Homepage (landing content)
			&homepage.HomepageHero{},
			&homepage.HomepageSection{},
			&homepage.HomepageTestimonial{},
		)

		if err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
		log.Println("Database migration completed successfully.")
	} else {
		log.Println("Database auto-migration is disabled (DB_AUTO_MIGRATE=false).")
	}

	return db, nil
}
