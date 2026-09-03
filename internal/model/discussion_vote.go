package model

import (
	"time"

	"github.com/google/uuid"
)
type DiscussionVote struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_vote_user_disc;not null"`
	DiscussionID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_vote_user_disc;index;not null"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
