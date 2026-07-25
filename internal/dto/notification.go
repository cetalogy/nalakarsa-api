package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID           uuid.UUID  `json:"id"`
	Type         string     `json:"type"`
	ActorID      *uuid.UUID `json:"actor_id,omitempty"`
	ActorName    string     `json:"actor_name,omitempty"`
	ResourceType string     `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty"`
	Message      string     `json:"message"`
	ReadAt       *time.Time `json:"read_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
