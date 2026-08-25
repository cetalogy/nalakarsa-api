package dto

import "github.com/google/uuid"

type InstitutionSuggestion struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	CountryCode     string    `json:"countryCode"`
	Country         string    `json:"country"`
	City            string    `json:"city"`
	Type            string    `json:"type"`
	IsInternational bool      `json:"isInternational"`
}
