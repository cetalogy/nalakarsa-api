package institutionrepository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nalakarsa/internal/model"
)

type InstitutionRepository interface {
	Search(search string, page, limit int) ([]model.Institution, int64, error)
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

func (r *institutionRepository) Search(search string, page, limit int) ([]model.Institution, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	search = strings.TrimSpace(search)

	query := r.db.Where("is_active = ?", true)
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	var total int64
	if err := query.Model(&model.Institution{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var result []model.Institution
	err := query.Order("name ASC").Offset((page - 1) * limit).Limit(limit).Find(&result).Error
	return result, total, err
}
