package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_proj_member;index;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_proj_member;index;not null"`
	Role      string    `gorm:"type:varchar(50);not null"`                  // role within the project
	Status    string    `gorm:"type:varchar(20);default:'active';not null"` // active, removed
	JoinedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
