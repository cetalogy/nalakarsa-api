package industryservice

import (
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/reference"
	industryrepository "nalakarsa/internal/repository/industry"
)

type IndustryService interface {
	Search(search string, limit int) ([]dto.IndustrySuggestion, error)
	Create(name string) (*dto.IndustrySuggestion, error)
}

func (s *industryService) Create(name string) (*dto.IndustrySuggestion, error) {
	name, err := reference.NormalizeName(name)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.Create(&model.Industry{Name: name, IsActive: true})
	if err != nil {
		return nil, err
	}
	return &dto.IndustrySuggestion{ID: item.ID, Name: item.Name}, nil
}

type industryService struct {
	repo industryrepository.IndustryRepository
}

func NewIndustryService(repo industryrepository.IndustryRepository) IndustryService {
	return &industryService{repo: repo}
}

func (s *industryService) Search(search string, limit int) ([]dto.IndustrySuggestion, error) {
	items, err := s.repo.Search(search, limit)
	if err != nil {
		return nil, err
	}
	result := make([]dto.IndustrySuggestion, len(items))
	for i, item := range items {
		result[i] = dto.IndustrySuggestion{ID: item.ID, Name: item.Name}
	}
	return result, nil
}
