package institutionservice

import (
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/reference"
	institutionrepository "nalakarsa/internal/repository/institution"
)

type InstitutionService interface {
	SearchInstitutions(search string, limit int) ([]dto.InstitutionSuggestion, error)
	CreateInstitution(name string) (*dto.InstitutionSuggestion, error)
}

func (s *institutionService) CreateInstitution(name string) (*dto.InstitutionSuggestion, error) {
	name, err := reference.NormalizeName(name)
	if err != nil {
		return nil, err
	}
	item, err := s.institutionRepo.Create(&model.Institution{Name: name, CountryCode: "ID", Country: "Indonesia", Type: "university", IsActive: true})
	if err != nil {
		return nil, err
	}
	return &dto.InstitutionSuggestion{ID: item.ID, Name: item.Name, CountryCode: item.CountryCode, Country: item.Country, City: item.City, Type: item.Type, IsInternational: strings.ToUpper(item.CountryCode) != "ID"}, nil
}

type institutionService struct {
	institutionRepo institutionrepository.InstitutionRepository
}

func NewInstitutionService(institutionRepo institutionrepository.InstitutionRepository) InstitutionService {
	return &institutionService{institutionRepo: institutionRepo}
}

func (s *institutionService) SearchInstitutions(search string, limit int) ([]dto.InstitutionSuggestion, error) {
	institutions, err := s.institutionRepo.Search(search, limit)
	if err != nil {
		return nil, err
	}

	res := make([]dto.InstitutionSuggestion, len(institutions))
	for i, institution := range institutions {
		isInternational := strings.ToUpper(strings.TrimSpace(institution.CountryCode)) != "ID"

		res[i] = dto.InstitutionSuggestion{
			ID:              institution.ID,
			Name:            institution.Name,
			CountryCode:     institution.CountryCode,
			Country:         institution.Country,
			City:            institution.City,
			Type:            institution.Type,
			IsInternational: isInternational,
		}
	}

	return res, nil
}
