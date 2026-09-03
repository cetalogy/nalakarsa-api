package institutionrepository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nalakarsa/internal/model"
)

type InstitutionRepository interface {
	Search(search string, limit int) ([]model.Institution, error)
	Create(item *model.Institution) (*model.Institution, error)
}

func (r *institutionRepository) Create(item *model.Institution) (*model.Institution, error) {
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
		return nil, err
	}
	var existing model.Institution
	if err := r.db.Where("name = ?", item.Name).First(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
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
