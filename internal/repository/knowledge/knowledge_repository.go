package knowledgerepository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nalakarsa/internal/model"
)

type KnowledgeRepository interface {
	SearchFields(search string, limit int) ([]model.KnowledgeField, error)
	Create(item *model.KnowledgeField) (*model.KnowledgeField, error)
}

func (r *knowledgeRepository) Create(item *model.KnowledgeField) (*model.KnowledgeField, error) {
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
		return nil, err
	}
	var existing model.KnowledgeField
	if err := r.db.Where("name = ?", item.Name).First(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

type knowledgeRepository struct{ db *gorm.DB }

func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository { return &knowledgeRepository{db: db} }

func normalizeLimit(limit int) int {
	if limit <= 0 || limit > 30 {
		return 10
	}
	return limit
}

func (r *knowledgeRepository) SearchFields(search string, limit int) ([]model.KnowledgeField, error) {
	var result []model.KnowledgeField
	q := r.db.Where("is_active = ?", true)
	if strings.TrimSpace(search) != "" {
		q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%")
	}
	err := q.Order("name ASC").Limit(normalizeLimit(limit)).Find(&result).Error
	return result, err
}
