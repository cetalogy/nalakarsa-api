package industryservice

import (
	"nalakarsa/internal/dto"
	industryrepository "nalakarsa/internal/repository/industry"
)

type IndustryService interface {
	Search(search string, limit int) ([]dto.IndustrySuggestion, error)
}

type industryService struct{ repo industryrepository.IndustryRepository }

func NewIndustryService(repo industryrepository.IndustryRepository) IndustryService {
	return &industryService{repo: repo}
}

func (s *industryService) Search(search string, limit int) ([]dto.IndustrySuggestion, error) {
	items, err := s.repo.Search(search, limit)
	if err != nil { return nil, err }
	result := make([]dto.IndustrySuggestion, len(items))
	for i, item := range items { result[i] = dto.IndustrySuggestion{ID: item.ID, Name: item.Name} }
	return result, nil
}
