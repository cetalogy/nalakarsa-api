package organizationservice

import (
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/reference"
	organizationrepository "nalakarsa/internal/repository/organization"
)

type OrganizationService interface {
	Search(search string, page, limit int) ([]dto.OrganizationSuggestion, int64, error)
	Create(name string) (*dto.OrganizationSuggestion, error)
}

type organizationService struct {
	repo organizationrepository.OrganizationRepository
}

func NewOrganizationService(repo organizationrepository.OrganizationRepository) OrganizationService {
	return &organizationService{repo: repo}
}

func (s *organizationService) Create(name string) (*dto.OrganizationSuggestion, error) {
	name, err := reference.NormalizeName(name)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.Create(&model.Organization{Name: name, CountryCode: "ID", Country: "Indonesia", IsActive: true})
	if err != nil {
		return nil, err
	}
	return toSuggestion(item), nil
}

func (s *organizationService) Search(search string, page, limit int) ([]dto.OrganizationSuggestion, int64, error) {
	items, total, err := s.repo.Search(search, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]dto.OrganizationSuggestion, len(items))
	for i := range items {
		result[i] = *toSuggestion(&items[i])
	}
	return result, total, nil
}

func toSuggestion(item *model.Organization) *dto.OrganizationSuggestion {
	return &dto.OrganizationSuggestion{ID: item.ID, Name: item.Name, CountryCode: strings.ToUpper(item.CountryCode), Country: item.Country}
}
