package expertiseservice

import (
	"nalakarsa/internal/dto"
	expertiserepository "nalakarsa/internal/repository/expertise"
)

type ExpertiseService interface {
	Search(search string, limit int) ([]dto.ExpertiseSuggestion, error)
}

type expertiseService struct {
	repo expertiserepository.ExpertiseRepository
}

func NewExpertiseService(repo expertiserepository.ExpertiseRepository) ExpertiseService {
	return &expertiseService{repo: repo}
}

func (s *expertiseService) Search(search string, limit int) ([]dto.ExpertiseSuggestion, error) {
	items, err := s.repo.Search(search, limit)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ExpertiseSuggestion, len(items))
	for i, item := range items {
		result[i] = dto.ExpertiseSuggestion{ID: item.ID, Name: item.Name, Category: item.Category, Specification: item.Specification}
	}
	return result, nil
}
