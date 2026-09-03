package expertiseservice

import (
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/reference"
	expertiserepository "nalakarsa/internal/repository/expertise"
)

type ExpertiseService interface {
	Search(search string, limit int) ([]dto.ExpertiseSuggestion, error)
	Create(name string) (*dto.ExpertiseSuggestion, error)
}

func (s *expertiseService) Create(name string) (*dto.ExpertiseSuggestion, error) {
	name, err := reference.NormalizeName(name); if err != nil { return nil, err }
	item, err := s.repo.Create(&model.Expertise{Name: name, IsActive: true}); if err != nil { return nil, err }
	return &dto.ExpertiseSuggestion{ID: item.ID, Name: item.Name}, nil
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
		result[i] = dto.ExpertiseSuggestion{ID: item.ID, Name: item.Name}
	}
	return result, nil
}
