package knowledgeservice

import (
	"nalakarsa/internal/dto"
	knowledgerepository "nalakarsa/internal/repository/knowledge"
)

type KnowledgeService interface {
	SearchFields(search string, limit int) ([]dto.KnowledgeFieldSuggestion, error)
}

type knowledgeService struct{ repo knowledgerepository.KnowledgeRepository }

func NewKnowledgeService(repo knowledgerepository.KnowledgeRepository) KnowledgeService { return &knowledgeService{repo: repo} }

func (s *knowledgeService) SearchFields(search string, limit int) ([]dto.KnowledgeFieldSuggestion, error) { items, err := s.repo.SearchFields(search, limit); if err != nil { return nil, err }; out := make([]dto.KnowledgeFieldSuggestion, len(items)); for i, item := range items { out[i] = dto.KnowledgeFieldSuggestion{ID:item.ID, Name:item.Name} }; return out, nil }
