package location

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"nalakarsa/internal/model"
)

//go:embed locations.json
var locationData []byte

type locationRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Province string `json:"province,omitempty"`
}

type LocationRepository interface {
	SearchLocations(search, locationType string, provinceID *uuid.UUID, limit int) ([]model.Location, error)
}

type locationRepository struct {
	locations []model.Location
	err       error
}

var (
	loadOnce sync.Once
	loaded   *locationRepository
)

func NewLocationRepository() LocationRepository {
	loadOnce.Do(func() {
		loaded = &locationRepository{}
		var records []locationRecord
		if err := json.Unmarshal(locationData, &records); err != nil {
			loaded.err = fmt.Errorf("load location data: %w", err)
			return
		}

		provinceIDs := make(map[string]uuid.UUID)
		for _, record := range records {
			if record.Type == "province" {
				provinceIDs[record.Name] = uuid.NewSHA1(uuid.Nil, []byte("province:"+record.Name))
			}
		}

		for _, record := range records {
			id := uuid.NewSHA1(uuid.Nil, []byte(record.Type+":"+record.Province+":"+record.Name))
			item := model.Location{ID: id, Name: record.Name, Type: record.Type, ProvinceName: record.Province, CountryCode: "ID", Country: "Indonesia", IsActive: true}
			if record.Type == "province" {
				item.ProvinceName = record.Name
			} else if provinceID, ok := provinceIDs[record.Province]; ok {
				item.ProvinceID = &provinceID
			}
			loaded.locations = append(loaded.locations, item)
		}
	})
	return loaded
}

func (r *locationRepository) SearchLocations(search, locationType string, provinceID *uuid.UUID, limit int) ([]model.Location, error) {
	if r.err != nil {
		return nil, r.err
	}
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	search = strings.ToLower(strings.TrimSpace(search))
	result := make([]model.Location, 0, limit)
	for _, item := range r.locations {
		if locationType != "" && item.Type != locationType {
			continue
		}
		if provinceID != nil && (item.ProvinceID == nil || *item.ProvinceID != *provinceID) {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(item.Name), search) {
			continue
		}
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type != result[j].Type {
			return result[i].Type < result[j].Type
		}
		return result[i].Name < result[j].Name
	})
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}
