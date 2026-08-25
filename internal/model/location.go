package model

import (
	"time"

	"github.com/google/uuid"
)

// Location stores normalized Indonesian province and city/master location data for autocomplete.
type Location struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Type        string    `gorm:"type:varchar(30);not null"` // province|city
	ProvinceID  *uuid.UUID `gorm:"type:uuid"`
	ProvinceName string   `gorm:"type:varchar(255);default:'';not null"`
	CountryCode string    `gorm:"type:varchar(10);default:'ID';not null"`
	Country     string    `gorm:"type:varchar(100);default:'Indonesia';not null"`
	IsActive    bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
