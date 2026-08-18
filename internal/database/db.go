package database

import (
	"fmt"
	"log"
	"os"
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
			&model.User{},
			&model.UserFollower{},

			&model.RefreshToken{},

			// Discussions
			&model.Discussion{},
			&model.DiscussionReply{},
			&model.DiscussionVote{},

			// Connections
			&model.Connection{},

			// Projects & Collaboration Requests
			&model.Project{},
			&model.ProjectMember{},
			&model.ProjectApplication{},
			&model.ProjectMilestone{},
			&model.CollaborationRequest{},

			// Chat & Group Chats
			&model.Conversation{},
			&model.ConversationMember{},
			&model.Message{},
			&model.GroupChat{},
			&model.GroupChatMember{},
			&model.GroupMessage{},

			// Notifications
			&model.Notification{},

			// Homepage (landing content)
			&model.HomepageHero{},
			&model.HomepageSection{},
			&model.HomepageTestimonial{},
		)

		if err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
		log.Println("Database migration completed successfully.")
	} else {
		log.Println("Database auto-migration is disabled (DB_AUTO_MIGRATE=false).")
	}

	// Ensure Counter Cache columns exist
	_ = db.Exec("ALTER TABLE discussions ADD COLUMN IF NOT EXISTS replies_count BIGINT DEFAULT 0 NOT NULL").Error
	_ = db.Exec("ALTER TABLE discussions ADD COLUMN IF NOT EXISTS upvote_count BIGINT DEFAULT 0 NOT NULL").Error

	// Ensure Compound Performance Indexes exist
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_replies_discussion_created ON discussion_replies (discussion_id, created_at DESC)").Error
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_discussion_user ON discussion_votes (discussion_id, user_id)").Error

	// Clean up legacy duplicated avatar URL prefixes in users table
	_ = db.Exec("UPDATE users SET avatar_url = REPLACE(avatar_url, '/avatars/avatars/avatars/', '/avatars/') WHERE avatar_url LIKE '%/avatars/avatars/avatars/%'").Error
	_ = db.Exec("UPDATE users SET avatar_url = REPLACE(avatar_url, '/avatars/avatars/', '/avatars/') WHERE avatar_url LIKE '%/avatars/avatars/%'").Error

	return db, nil
}
