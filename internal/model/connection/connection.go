package connection

import (
	"time"

	"nalakarsa/internal/model/user"

	"github.com/google/uuid"
)

// Connection represents a network connection between two users.
// Status: pending, accepted, rejected
// Unique constraint prevents duplicate connection requests.
type Connection struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RequesterID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_conn_pair;index;not null"`
	AddresseeID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_conn_pair;index;not null"`
	Status      string    `gorm:"type:varchar(20);index;default:'pending';not null"` // pending, accepted, rejected
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Requester user.User `gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	Addressee user.User `gorm:"foreignKey:AddresseeID;constraint:OnDelete:CASCADE"`
}
