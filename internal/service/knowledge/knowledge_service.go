package knowledgeservice

import (
	"nalakarsa/internal/dto"
	knowledgerepository "nalakarsa/internal/repository/knowledge"

	"github.com/google/uuid"
)

type KnowledgeService interface {
	SearchFields(subdomainID uuid.UUID, search string, limit int) ([]dto.KnowledgeFieldSuggestion, error)
}

type knowledgeService struct{ repo knowledgerepository.KnowledgeRepository }

func NewKnowledgeService(repo knowledgerepository.KnowledgeRepository) KnowledgeService { return &knowledgeService{repo: repo} }

func (s *knowledgeService) SearchFields(id uuid.UUID, search string, limit int) ([]dto.KnowledgeFieldSuggestion, error) { subdomainID := ""; if id != uuid.Nil { subdomainID = id.String() }; items, err := s.repo.SearchFields(subdomainID, search, limit); if err != nil { return nil, err }; out := make([]dto.KnowledgeFieldSuggestion, len(items)); for i, item := range items { out[i] = dto.KnowledgeFieldSuggestion{ID:item.ID, SubdomainID:item.SubdomainID, SubdomainName:item.SubdomainName, DomainID:item.DomainID, DomainName:item.DomainName, Name:item.Name} }; return out, nil }
