package dto

import "github.com/google/uuid"

type KnowledgeFieldSuggestion struct {
	ID            uuid.UUID `json:"id"`
	SubdomainID   uuid.UUID `json:"subdomain_id"`
	SubdomainName string    `json:"subdomain_name"`
	DomainID      uuid.UUID `json:"domain_id"`
	DomainName    string    `json:"domain_name"`
	Name          string    `json:"name"`
}
