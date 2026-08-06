package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID           uuid.UUID  `json:"id" form:"id"`
	Type         string     `json:"type" form:"type"`
	ActorID      *uuid.UUID `json:"actor_id,omitempty" form:"actor_id"`
	ActorName    string     `json:"actor_name,omitempty" form:"actor_name"`
	ResourceType string     `json:"resource_type,omitempty" form:"resource_type"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty" form:"resource_id"`
	Message      string     `json:"message" form:"message"`
	ReadAt       *time.Time `json:"read_at" form:"read_at"`
	CreatedAt    time.Time  `json:"created_at" form:"created_at"`
}
