package expertiserepository

import (
	"strings"

	"nalakarsa/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExpertiseRepository interface {
	Search(search string, limit int) ([]model.Expertise, error)
	Create(item *model.Expertise) (*model.Expertise, error)
}

func (r *expertiseRepository) Create(item *model.Expertise) (*model.Expertise, error) {
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
		return nil, err
	}
	var existing model.Expertise
	if err := r.db.Where("name = ?", item.Name).First(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

type expertiseRepository struct {
	db *gorm.DB
}

func NewExpertiseRepository(db *gorm.DB) ExpertiseRepository {
	return &expertiseRepository{db: db}
}

func (r *expertiseRepository) Search(search string, limit int) ([]model.Expertise, error) {
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	query := r.db.Where("is_active = ?", true)
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ?", like)
	}
	var result []model.Expertise
	err := query.Order("name ASC").Limit(limit).Find(&result).Error
	return result, err
}
