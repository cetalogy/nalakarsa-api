package dto

import "github.com/google/uuid"

type SkillResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateSkillRequest struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}

type AddUserSkillRequest struct {
	SkillID *uuid.UUID `json:"skill_id"`
	Name    *string    `json:"name"`
}
