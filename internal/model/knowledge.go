package model

import (
	"time"

	"github.com/google/uuid"
)

type KnowledgeField struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsActive  bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
