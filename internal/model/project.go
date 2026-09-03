package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Project struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OwnerID            uuid.UUID      `gorm:"type:uuid;index;not null"`
	Title              string         `gorm:"type:varchar(255);not null"`
	Description        string         `gorm:"type:text;not null"`
	Category           string         `gorm:"type:varchar(100);index;not null"`
	Status             string         `gorm:"type:varchar(20);index;default:'draft';not null"` // draft, open, in_review, active, completed, archived
	FundingStatus      string         `gorm:"type:varchar(50);default:'';not null"`
	Location           string         `gorm:"type:varchar(255);default:'';not null"`
	Needs              string         `gorm:"type:varchar(50);default:'';not null"` // akademisi, praktisi, profesional
	Deadline           *time.Time     `gorm:"type:timestamptz"`
	Progress           int            `gorm:"type:int;default:0;not null"` // 0-100
	SourceDiscussionID *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt          time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	Owner        User                 `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Members      []ProjectMember      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Applications []ProjectApplication `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Milestones   []ProjectMilestone   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}
