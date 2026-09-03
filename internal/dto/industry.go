package dto

import "github.com/google/uuid"

type IndustrySuggestion struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
