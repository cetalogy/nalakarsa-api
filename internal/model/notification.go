package model

import (
	"time"

	"github.com/google/uuid"
)
type Notification struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID  `gorm:"type:uuid;index;not null"` // recipient
	Type         string     `gorm:"type:varchar(50);index;not null"`
	ActorID      *uuid.UUID `gorm:"type:uuid;index"`                      // who triggered the notification
	ResourceType string     `gorm:"type:varchar(50);default:'';not null"` // discussion, project, connection, etc.
	ResourceID   *uuid.UUID `gorm:"type:uuid"`
	Payload      string     `gorm:"type:jsonb;default:'{}';not null"` // additional JSON data
	ReadAt       *time.Time `gorm:"type:timestamptz"`
	CreatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Actor *User `gorm:"foreignKey:ActorID;constraint:OnDelete:SET NULL"`
}
