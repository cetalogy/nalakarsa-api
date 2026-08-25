package institutionrepository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"nalakarsa/internal/model"
)

type InstitutionRepository interface {
	Search(search string, limit int) ([]model.Institution, error)
}

type hipolabsInstitution struct {
	Name        string `json:"name"`
	Country     string `json:"country"`
	CountryCode string `json:"alpha_two_code"`
	State       string `json:"state-province"`
}

type apiInstitutionRepository struct {
	client  *http.Client
	baseURL string
}

func NewInstitutionRepository() InstitutionRepository {
	return &apiInstitutionRepository{client: &http.Client{Timeout: 8 * time.Second}, baseURL: "https://universities.hipolabs.com/search"}
}

func (r *apiInstitutionRepository) Search(search string, limit int) ([]model.Institution, error) {
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	search = strings.TrimSpace(search)
	if search == "" {
		return []model.Institution{}, nil
	}

	endpoint, err := url.Parse(r.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse university API URL: %w", err)
	}
	query := endpoint.Query()
	query.Set("name", search)
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create university API request: %w", err)
	}
	response, err := r.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request university API: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("university API returned status %s", response.Status)
	}

	var payload []hipolabsInstitution
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode university API response: %w", err)
	}

	result := make([]model.Institution, 0, len(payload))
	for _, item := range payload {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		result = append(result, model.Institution{
			ID:          uuid.NewSHA1(uuid.Nil, []byte(item.CountryCode+":"+item.Name)),
			Name:        item.Name,
			CountryCode: strings.ToUpper(item.CountryCode),
			Country:     item.Country,
			City:        item.State,
			Type:        "university",
			IsActive:    true,
		})
	}

	sort.SliceStable(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}
