package institutionrepository

import (
	"strings"

	"nalakarsa/internal/model"
	"gorm.io/gorm"
)

type InstitutionRepository interface {
	Search(search string, limit int) ([]model.Institution, error)
}

type institutionRepository struct{ db *gorm.DB }

func NewInstitutionRepository(db *gorm.DB) InstitutionRepository {
	return &institutionRepository{db: db}
}

func (r *institutionRepository) Search(search string, limit int) ([]model.Institution, error) {
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	search = strings.TrimSpace(search)
	if search == "" {
		return []model.Institution{}, nil
	}

	var result []model.Institution
	err := r.db.Where("is_active = ? AND name ILIKE ?", true, "%"+search+"%").
		Order("name ASC").Limit(limit).Find(&result).Error
	return result, err
}
