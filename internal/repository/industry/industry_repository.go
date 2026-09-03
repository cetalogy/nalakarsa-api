package industryrepository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nalakarsa/internal/model"
)

type IndustryRepository interface {
	Search(search string, limit int) ([]model.Industry, error)
	Create(item *model.Industry) (*model.Industry, error)
}

func (r *industryRepository) Create(item *model.Industry) (*model.Industry, error) {
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
		return nil, err
	}
	var existing model.Industry
	if err := r.db.Where("name = ?", item.Name).First(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

type industryRepository struct{ db *gorm.DB }

func NewIndustryRepository(db *gorm.DB) IndustryRepository { return &industryRepository{db: db} }

func (r *industryRepository) Search(search string, limit int) ([]model.Industry, error) {
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	query := r.db.Where("is_active = ?", true)
	if search = strings.TrimSpace(search); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	var result []model.Industry
	err := query.Order("name ASC").Limit(limit).Find(&result).Error
	return result, err
}
