package institutionservice

import (
	"strings"

	"nalakarsa/internal/dto"
	institutionrepository "nalakarsa/internal/repository/institution"
)

type InstitutionService interface {
	SearchInstitutions(search string, limit int) ([]dto.InstitutionSuggestion, error)
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
