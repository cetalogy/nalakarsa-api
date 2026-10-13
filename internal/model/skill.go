package model

import (
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type UserSkill struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_skill;index"`
	SkillID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_skill;index"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Skill     Skill     `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}
