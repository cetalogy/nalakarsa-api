package model

import (
	"time"

	"github.com/google/uuid"
)
type Institution struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	CountryCode string    `gorm:"type:varchar(10);default:'';not null"`
	Country     string    `gorm:"type:varchar(100);default:'';not null"`
	City        string    `gorm:"type:varchar(100);default:'';not null"`
	Type        string    `gorm:"type:varchar(50);default:'university';not null"`
	IsActive    bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
