package industryrepository

import (
	"strings"

	"nalakarsa/internal/model"
	"gorm.io/gorm"
)

type IndustryRepository interface {
	Search(search string, limit int) ([]model.Industry, error)
}

type industryRepository struct{ db *gorm.DB }

func NewIndustryRepository(db *gorm.DB) IndustryRepository { return &industryRepository{db: db} }

func (r *industryRepository) Search(search string, limit int) ([]model.Industry, error) {
	if limit <= 0 || limit > 30 { limit = 10 }
	query := r.db.Where("is_active = ?", true)
	if search = strings.TrimSpace(search); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	var result []model.Industry
	err := query.Order("name ASC").Limit(limit).Find(&result).Error
	return result, err
}
