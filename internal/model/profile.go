package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID         uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"`
	NamaLengkap    string         `gorm:"type:varchar(255);not null"`
	GelarDepan     string         `gorm:"type:varchar(100);default:'';not null"`
	GelarBelakang  string         `gorm:"type:varchar(100);default:'';not null"`
	Afiliasi       string         `gorm:"type:varchar(255);not null"`
	Lokasi         string         `gorm:"type:varchar(255);not null"`
	BidangKeahlian string         `gorm:"type:varchar(255);not null"`
	Industry       string         `gorm:"type:varchar(255);default:'';not null"`
	Bio            string         `gorm:"type:text;default:'';not null"`
	Mission        string         `gorm:"type:text;default:'';not null"`
	AvatarURL      string         `gorm:"type:varchar(255);default:'';not null"`
	ViewCount      int            `gorm:"type:int;default:0;not null"`
	CreatedAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
