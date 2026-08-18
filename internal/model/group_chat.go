package model

import (
	"time"

	"github.com/google/uuid"
)

// GroupChat represents a collaboration project group chat room.
type GroupChat struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TopicID         *uuid.UUID  `gorm:"type:uuid;index"`
	ProjectID       *uuid.UUID  `gorm:"type:uuid;index"`
	Title           string      `gorm:"type:varchar(255);not null"`
	Badge           string      `gorm:"type:varchar(50);default:'Grup Kolaborasi'"`
	LastMessage     *string     `gorm:"type:text"`
	LastMessageTime *time.Time  `gorm:"type:timestamptz"`
	CreatedAt       time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Topic    *Discussion       `gorm:"foreignKey:TopicID;constraint:OnDelete:SET NULL"`
	Project  *Project          `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	Members  []GroupChatMember `gorm:"foreignKey:GroupChatID;constraint:OnDelete:CASCADE"`
	Messages []GroupMessage    `gorm:"foreignKey:GroupChatID;constraint:OnDelete:CASCADE"`
}

// GroupChatMember represents a member in a collaboration group chat.
type GroupChatMember struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	GroupChatID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_group_chat_member,priority:1;index;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_group_chat_member,priority:2;index;not null"`
	Role        string    `gorm:"type:varchar(50);default:'Mitra Kolaborasi'"`
	JoinedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// GroupMessage represents a message sent within a collaboration group chat.
type GroupMessage struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	GroupChatID     uuid.UUID  `gorm:"type:uuid;index;not null"`
	SenderID        *uuid.UUID `gorm:"type:uuid;index"` // NULL if system message
	IsSystemMessage bool       `gorm:"default:false;not null"`
	Content         string     `gorm:"type:text;not null"`
	CreatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Sender *User `gorm:"foreignKey:SenderID;constraint:OnDelete:SET NULL"`
}
