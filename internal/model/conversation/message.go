package conversation

import (
	"time"

	"nalakarsa/internal/model/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message represents a single chat message in a conversation.
// Status: sent, delivered, read
type Message struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ConversationID uuid.UUID      `gorm:"type:uuid;index;not null"`
	SenderID       uuid.UUID      `gorm:"type:uuid;index;not null"`
	Body           string         `gorm:"type:text;not null"`
	Status         string         `gorm:"type:varchar(20);default:'sent';not null"` // sent, delivered, read
	CreatedAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Relations
	Sender user.User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
}
