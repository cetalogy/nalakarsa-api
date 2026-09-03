package knowledgerepository

import (
	"strings"

	"nalakarsa/internal/model"
	"gorm.io/gorm"
)

type KnowledgeRepository interface {
	SearchFields(search string, limit int) ([]model.KnowledgeField, error)
}

type knowledgeRepository struct{ db *gorm.DB }

func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository { return &knowledgeRepository{db: db} }

func normalizeLimit(limit int) int { if limit <= 0 || limit > 30 { return 10 }; return limit }

func (r *knowledgeRepository) SearchFields(search string, limit int) ([]model.KnowledgeField, error) {
	var result []model.KnowledgeField
	q := r.db.Where("is_active = ?", true)
	if strings.TrimSpace(search) != "" { q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%") }
	err := q.Order("name ASC").Limit(normalizeLimit(limit)).Find(&result).Error
	return result, err
}
