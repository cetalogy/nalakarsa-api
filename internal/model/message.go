package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ConversationID     uuid.UUID      `gorm:"type:uuid;index;not null"`
	SenderID           uuid.UUID      `gorm:"type:uuid;index;not null"`
	Body               string         `gorm:"type:text;not null"`
	Status             string         `gorm:"type:varchar(20);default:'sent';not null"`
	AttachmentPath     string         `gorm:"type:text"`
	AttachmentName     string         `gorm:"type:varchar(255)"`
	AttachmentMimeType string         `gorm:"type:varchar(100)"`
	AttachmentSize     int64          `gorm:"default:0"`
	AttachmentType     string         `gorm:"type:varchar(20)"`
	CreatedAt          time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	Sender User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
}
