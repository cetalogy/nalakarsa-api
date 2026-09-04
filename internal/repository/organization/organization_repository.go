package organizationrepository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nalakarsa/internal/model"
)

type OrganizationRepository interface {
	Search(search string, page, limit int) ([]model.Organization, int64, error)
	Create(item *model.Organization) (*model.Organization, error)
}

type organizationRepository struct{ db *gorm.DB }

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(item *model.Organization) (*model.Organization, error) {
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
		return nil, err
	}
	var existing model.Organization
	if err := r.db.Where("name = ?", item.Name).First(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *organizationRepository) Search(search string, page, limit int) ([]model.Organization, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	query := r.db.Where("is_active = ?", true)
	if search = strings.TrimSpace(search); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	var total int64
	if err := query.Model(&model.Organization{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var result []model.Organization
	err := query.Order("name ASC").Offset((page - 1) * limit).Limit(limit).Find(&result).Error
	return result, total, err
}
