package model

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	CountryCode string    `gorm:"type:varchar(10);not null;default:'ID'"`
	Country     string    `gorm:"type:varchar(100);not null;default:'Indonesia'"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
