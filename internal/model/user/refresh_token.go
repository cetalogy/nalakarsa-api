package user

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	Token      string    `gorm:"type:varchar(512);uniqueIndex;not null"`
	DeviceInfo string    `gorm:"type:varchar(255);default:'';not null"`
	UserAgent  string    `gorm:"type:varchar(255);default:'';not null"`
	IPAddress  string    `gorm:"type:varchar(64);default:'';not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
