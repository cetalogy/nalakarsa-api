package conversation

import (
	"time"

	"github.com/google/uuid"
)

// Conversation represents a chat conversation.
// Type: "direct" for one-on-one conversations.
type Conversation struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Type          string     `gorm:"type:varchar(20);default:'direct';not null"` // direct
	LastMessageAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Members  []ConversationMember `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
	Messages []Message            `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
}
