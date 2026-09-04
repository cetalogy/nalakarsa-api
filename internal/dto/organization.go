package dto

import "github.com/google/uuid"

type OrganizationSuggestion struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CountryCode string    `json:"countryCode"`
	Country     string    `json:"country"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}
