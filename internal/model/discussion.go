package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Discussion struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID             uuid.UUID      `gorm:"type:uuid;index;not null"`
	Title              string         `gorm:"type:varchar(255);not null"`
	Description        string         `gorm:"type:text;not null"`
	Category           string         `gorm:"type:varchar(100);index;not null"`
	Tags               string         `gorm:"type:varchar(255);default:'';not null"`
	Status             string         `gorm:"type:varchar(20);index;default:'open';not null"` // open, resolved, closed
	IsInCollaboration  bool           `gorm:"default:false;not null"`
	SourceDiscussionID *uuid.UUID     `gorm:"type:uuid;index"`
	RepliesCount       int64          `gorm:"default:0;not null"`
	UpvoteCount        int64          `gorm:"default:0;not null"`
	CreatedAt          time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	// Relations
	User    User              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Replies []DiscussionReply `gorm:"foreignKey:DiscussionID;constraint:OnDelete:CASCADE"`
	Votes   []DiscussionVote  `gorm:"foreignKey:DiscussionID;constraint:OnDelete:CASCADE"`
}
