package conversation

import (
	"time"

	"nalakarsa/internal/model/user"

	"github.com/google/uuid"
)

type ConversationMember struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ConversationID    uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_conv_member;index;not null"`
	UserID            uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_conv_member;index;not null"`
	LastReadMessageID *uuid.UUID `gorm:"type:uuid"`
	JoinedAt          time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	User user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
