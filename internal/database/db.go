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
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	if cfg.Env == "production" {
		dbLogger = logger.New(
			log.New(os.Stdout, "", 0),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
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
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection successfully established.")
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}
	if cfg.DBAutoMigrate {
		log.Println("Running database migrations...")
		err = db.AutoMigrate(
			&model.User{},
			&model.UserFollower{},
			&model.Institution{},
			&model.Expertise{},
			&model.Industry{},
			&model.KnowledgeField{},

			&model.RefreshToken{},
			&model.Discussion{},
			&model.DiscussionReply{},
			&model.DiscussionVote{},
			&model.Connection{},
			&model.Project{},
			&model.ProjectMember{},
			&model.ProjectApplication{},
			&model.ProjectMilestone{},
			&model.CollaborationRequest{},
			&model.Conversation{},
			&model.ConversationMember{},
			&model.Message{},
			&model.GroupChat{},
			&model.GroupChatMember{},
			&model.GroupMessage{},
			&model.Notification{},
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
	latestSchemaStatements := []string{
		`DROP INDEX IF EXISTS uni_expertises_name`,
		`DROP INDEX IF EXISTS idx_expertises_name`,
		`DROP INDEX IF EXISTS idx_expertise_hierarchy`,
		`DO $$ BEGIN IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'expertises' AND column_name = 'specification') THEN UPDATE expertises SET name = trim(specification) WHERE trim(specification) <> ''; END IF; END $$`,
		`DELETE FROM expertises WHERE trim(name) = ''`,
		`DELETE FROM expertises a USING expertises b WHERE a.id > b.id AND a.name = b.name`,
		`ALTER TABLE expertises DROP COLUMN IF EXISTS category`,
		`ALTER TABLE expertises DROP COLUMN IF EXISTS specification`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_expertises_name ON expertises (name)`,
		`ALTER TABLE knowledge_fields DROP CONSTRAINT IF EXISTS fk_knowledge_fields_subdomain`,
		`DROP INDEX IF EXISTS idx_knowledge_field_name`,
		`ALTER TABLE knowledge_fields DROP COLUMN IF EXISTS subdomain_id`,
		`DELETE FROM knowledge_fields a USING knowledge_fields b WHERE a.id > b.id AND a.name = b.name`,
		`DROP TABLE IF EXISTS knowledge_subdomains`,
		`DROP TABLE IF EXISTS knowledge_domains`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_knowledge_field_name ON knowledge_fields (name)`,
		`INSERT INTO industries (name, is_active) VALUES
			('Pertanian, Kehutanan dan Perikanan', true),
			('Pertambangan dan Penggalian', true),
			('Industri Pengolahan', true),
			('Pengadaan Listrik, Gas, Uap/Air Panas Dan Udara Dingin', true),
			('Treatment Air, Limbah, Pemulihan Material Sampah & Remediasi', true),
			('Konstruksi', true),
			('Perdagangan Besar & Eceran; Reparasi Mobil & Sepeda Motor', true),
			('Pengangkutan dan Pergudangan', true),
			('Penyediaan Akomodasi Dan Penyediaan Makan Minum', true),
			('Informasi Dan Komunikasi', true),
			('Aktivitas Keuangan dan Asuransi', true),
			('Real Estat', true),
			('Aktivitas Profesional, Ilmiah Dan Teknis', true),
			('Aktivitas Penyewaan, Ketenagakerjaan, Agen Perjalanan & Penunjang', true),
			('Administrasi Pemerintahan, Pertahanan & Jaminan Sosial', true),
			('Pendidikan', true),
			('Aktivitas Kesehatan Manusia Dan Aktivitas Sosial', true),
			('Kesenian, Hiburan Dan Rekreasi', true),
			('Aktivitas Jasa Lainnya', true),
			('Aktivitas Rumah Tangga Sebagai Pemberi Kerja', true),
			('Aktivitas Badan Internasional & Badan Ekstra Internasional', true)
		ON CONFLICT (name) DO NOTHING`,
		`DROP INDEX IF EXISTS idx_knowledge_field_name`,
		`UPDATE knowledge_fields SET name = trim(regexp_replace(regexp_replace(name, '^[[:space:]]*[(]?[0-9]+[.)][[:space:]]*', ''), '[[:space:]]+[0-9]+$', '')) WHERE name ~ '^[[:space:]]*[(]?[0-9]+[.)]' OR name ~ '[[:space:]]+[0-9]+$'`,
		`DELETE FROM knowledge_fields a USING knowledge_fields b WHERE a.id > b.id AND a.name = b.name`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_knowledge_field_name ON knowledge_fields (name)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS security_question VARCHAR(30) NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS security_answer_hash VARCHAR(255) NOT NULL DEFAULT ''`,
		`ALTER TABLE messages ADD COLUMN IF NOT EXISTS attachment_path TEXT`,
		`ALTER TABLE messages ADD COLUMN IF NOT EXISTS attachment_name VARCHAR(255)`,
		`ALTER TABLE messages ADD COLUMN IF NOT EXISTS attachment_mime_type VARCHAR(100)`,
		`ALTER TABLE messages ADD COLUMN IF NOT EXISTS attachment_size BIGINT DEFAULT 0`,
		`ALTER TABLE messages ADD COLUMN IF NOT EXISTS attachment_type VARCHAR(20)`,
		`ALTER TABLE group_messages ADD COLUMN IF NOT EXISTS attachment_path TEXT`,
		`ALTER TABLE group_messages ADD COLUMN IF NOT EXISTS attachment_name VARCHAR(255)`,
		`ALTER TABLE group_messages ADD COLUMN IF NOT EXISTS attachment_mime_type VARCHAR(100)`,
		`ALTER TABLE group_messages ADD COLUMN IF NOT EXISTS attachment_size BIGINT DEFAULT 0`,
		`ALTER TABLE group_messages ADD COLUMN IF NOT EXISTS attachment_type VARCHAR(20)`,
	}
	for _, statement := range latestSchemaStatements {
		if err := db.Exec(statement).Error; err != nil {
			return nil, fmt.Errorf("failed to apply latest schema change: %w", err)
		}
	}
	_ = db.Exec("ALTER TABLE discussions ADD COLUMN IF NOT EXISTS replies_count BIGINT DEFAULT 0 NOT NULL").Error
	_ = db.Exec("ALTER TABLE discussions ADD COLUMN IF NOT EXISTS upvote_count BIGINT DEFAULT 0 NOT NULL").Error
	_ = db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error
	_ = db.Exec("DELETE FROM refresh_tokens WHERE expires_at <= NOW()").Error
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens (expires_at)").Error
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications (user_id, created_at DESC)").Error
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages (conversation_id, created_at DESC)").Error
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_connections_requester_status ON connections (requester_id, status)").Error
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_connections_addressee_status ON connections (addressee_id, status)").Error
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_replies_discussion_created ON discussion_replies (discussion_id, created_at DESC)").Error
	_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_discussion_user ON discussion_votes (discussion_id, user_id)").Error
	_ = db.Exec("UPDATE users SET avatar_url = REPLACE(avatar_url, '/avatars/avatars/avatars/', '/avatars/') WHERE avatar_url LIKE '%/avatars/avatars/avatars/%'").Error
	_ = db.Exec("UPDATE users SET avatar_url = REPLACE(avatar_url, '/avatars/avatars/', '/avatars/') WHERE avatar_url LIKE '%/avatars/avatars/%'").Error

	return db, nil
}
func StartRefreshTokenCleanup(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := db.Exec("DELETE FROM refresh_tokens WHERE expires_at <= NOW()").Error; err != nil {
				log.Printf("refresh token cleanup failed: %v", err)
			}
		}
	}()
}
