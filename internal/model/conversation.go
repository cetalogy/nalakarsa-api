package model

import (
	"time"

	"github.com/google/uuid"
)
type Conversation struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Type          string     `gorm:"type:varchar(20);default:'direct';not null"` // direct
	LastMessageAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Members  []ConversationMember `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
	Messages []Message            `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE"`
}
