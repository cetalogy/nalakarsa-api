package dto

import "github.com/google/uuid"

type KnowledgeFieldSuggestion struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
