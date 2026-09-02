package knowledgeservice

import (
	"nalakarsa/internal/dto"
	knowledgerepository "nalakarsa/internal/repository/knowledge"

	"github.com/google/uuid"
)

type KnowledgeService interface {
	SearchDomains(search string, limit int) ([]dto.KnowledgeDomainSuggestion, error)
	SearchSubdomains(domainID uuid.UUID, search string, limit int) ([]dto.KnowledgeSubdomainSuggestion, error)
	SearchFields(subdomainID uuid.UUID, search string, limit int) ([]dto.KnowledgeFieldSuggestion, error)
}

type knowledgeService struct{ repo knowledgerepository.KnowledgeRepository }

func NewKnowledgeService(repo knowledgerepository.KnowledgeRepository) KnowledgeService { return &knowledgeService{repo: repo} }

func (s *knowledgeService) SearchDomains(search string, limit int) ([]dto.KnowledgeDomainSuggestion, error) { items, err := s.repo.SearchDomains(search, limit); if err != nil { return nil, err }; out := make([]dto.KnowledgeDomainSuggestion, len(items)); for i, item := range items { out[i] = dto.KnowledgeDomainSuggestion{ID:item.ID, Name:item.Name} }; return out, nil }
func (s *knowledgeService) SearchSubdomains(id uuid.UUID, search string, limit int) ([]dto.KnowledgeSubdomainSuggestion, error) { items, err := s.repo.SearchSubdomains(id.String(), search, limit); if err != nil { return nil, err }; out := make([]dto.KnowledgeSubdomainSuggestion, len(items)); for i, item := range items { out[i] = dto.KnowledgeSubdomainSuggestion{ID:item.ID, DomainID:item.DomainID, Name:item.Name} }; return out, nil }
func (s *knowledgeService) SearchFields(id uuid.UUID, search string, limit int) ([]dto.KnowledgeFieldSuggestion, error) { items, err := s.repo.SearchFields(id.String(), search, limit); if err != nil { return nil, err }; out := make([]dto.KnowledgeFieldSuggestion, len(items)); for i, item := range items { out[i] = dto.KnowledgeFieldSuggestion{ID:item.ID, SubdomainID:item.SubdomainID, Name:item.Name} }; return out, nil }
