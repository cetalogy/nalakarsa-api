package expertiserepository

import (
	"strings"

	"nalakarsa/internal/model"

	"gorm.io/gorm"
)

type ExpertiseRepository interface {
	Search(search string, limit int) ([]model.Expertise, error)
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
		query = query.Where("name ILIKE ? OR category ILIKE ? OR specification ILIKE ?", like, like, like)
	}
	var result []model.Expertise
	err := query.Order("name ASC").Limit(limit).Find(&result).Error
	return result, err
}
