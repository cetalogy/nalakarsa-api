package dto

import "github.com/google/uuid"

type ExpertiseSuggestion struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
