package model

import (
	"time"

	"github.com/google/uuid"
)

type Expertise struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name          string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_expertise_hierarchy"`
	Category      string    `gorm:"type:varchar(100);default:'';not null;uniqueIndex:idx_expertise_hierarchy"`
	Specification string    `gorm:"type:varchar(255);default:'';not null;uniqueIndex:idx_expertise_hierarchy"`
	IsActive  bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
