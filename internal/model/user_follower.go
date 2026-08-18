package model

import (
	"time"

	"github.com/google/uuid"
)

// UserFollower represents a one-way follow relationship between users.
// FollowerID is the user who follows; FollowingID is the user being followed.
type UserFollower struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FollowerID  uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_followers_pair,priority:1;index;not null"`
	FollowingID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_followers_pair,priority:2;index;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Follower  User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	Following User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE"`
}
