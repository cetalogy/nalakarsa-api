package locationservice

import (
	"fmt"
	"strings"

	"nalakarsa/internal/dto"
	locationrepository "nalakarsa/internal/repository/location"

	"github.com/google/uuid"
)

type LocationService interface {
	SearchLocations(search, locationType string, provinceID string, limit int) ([]dto.LocationSuggestion, error)
}

type locationService struct {
	locationRepo locationrepository.LocationRepository
}

func NewLocationService(locationRepo locationrepository.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

func (s *locationService) SearchLocations(search, locationType string, provinceID string, limit int) ([]dto.LocationSuggestion, error) {
	var parsedProvinceID *uuid.UUID
	if strings.TrimSpace(provinceID) != "" {
		id, err := uuid.Parse(provinceID)
		if err != nil {
			return nil, fmt.Errorf("invalid provinceId: %w", err)
		}
		parsedProvinceID = &id
	}

	locations, err := s.locationRepo.SearchLocations(search, strings.TrimSpace(locationType), parsedProvinceID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]dto.LocationSuggestion, 0, len(locations))
	for _, loc := range locations {
		result = append(result, dto.LocationSuggestion{
			ID:           loc.ID,
			Name:         loc.Name,
			Type:         loc.Type,
			ProvinceID:   loc.ProvinceID,
			ProvinceName: loc.ProvinceName,
			IsProvince:   strings.EqualFold(loc.Type, "province"),
			IsCity:       strings.EqualFold(loc.Type, "city"),
		})
	}

	return result, nil
}
