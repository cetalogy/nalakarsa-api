package knowledgeservice

import (
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/reference"
	knowledgerepository "nalakarsa/internal/repository/knowledge"
)

type KnowledgeService interface {
	SearchFields(search string, limit int) ([]dto.KnowledgeFieldSuggestion, error)
	CreateField(name string) (*dto.KnowledgeFieldSuggestion, error)
}

func (s *knowledgeService) CreateField(name string) (*dto.KnowledgeFieldSuggestion, error) {
	name, err := reference.NormalizeName(name)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.Create(&model.KnowledgeField{Name: name, IsActive: true})
	if err != nil {
		return nil, err
	}
	return &dto.KnowledgeFieldSuggestion{ID: item.ID, Name: item.Name}, nil
}

type knowledgeService struct {
	repo knowledgerepository.KnowledgeRepository
}

func NewKnowledgeService(repo knowledgerepository.KnowledgeRepository) KnowledgeService {
	return &knowledgeService{repo: repo}
}

func (s *knowledgeService) SearchFields(search string, limit int) ([]dto.KnowledgeFieldSuggestion, error) {
	items, err := s.repo.SearchFields(search, limit)
	if err != nil {
		return nil, err
	}
	out := make([]dto.KnowledgeFieldSuggestion, len(items))
	for i, item := range items {
		out[i] = dto.KnowledgeFieldSuggestion{ID: item.ID, Name: item.Name}
	}
	return out, nil
}
