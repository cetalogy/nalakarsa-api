package knowledgerepository

import (
	"strings"

	"nalakarsa/internal/model"
	"gorm.io/gorm"
)

type KnowledgeRepository interface {
	SearchDomains(search string, limit int) ([]model.KnowledgeDomain, error)
	SearchSubdomains(domainID, search string, limit int) ([]model.KnowledgeSubdomain, error)
	SearchFields(subdomainID, search string, limit int) ([]model.KnowledgeField, error)
}

type knowledgeRepository struct{ db *gorm.DB }

func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository { return &knowledgeRepository{db: db} }

func normalizeLimit(limit int) int { if limit <= 0 || limit > 30 { return 10 }; return limit }

func (r *knowledgeRepository) SearchDomains(search string, limit int) ([]model.KnowledgeDomain, error) {
	var result []model.KnowledgeDomain
	q := r.db.Where("is_active = ?", true); if strings.TrimSpace(search) != "" { q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%") }
	err := q.Order("name ASC").Limit(normalizeLimit(limit)).Find(&result).Error; return result, err
}

func (r *knowledgeRepository) SearchSubdomains(domainID, search string, limit int) ([]model.KnowledgeSubdomain, error) {
	var result []model.KnowledgeSubdomain
	q := r.db.Where("is_active = ? AND domain_id = ?", true, domainID); if strings.TrimSpace(search) != "" { q = q.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%") }
	err := q.Order("name ASC").Limit(normalizeLimit(limit)).Find(&result).Error; return result, err
}

func (r *knowledgeRepository) SearchFields(subdomainID, search string, limit int) ([]model.KnowledgeField, error) {
	var result []model.KnowledgeField
	q := r.db.Table("knowledge_fields AS f").
		Select("f.*, s.name AS subdomain_name, s.domain_id, d.name AS domain_name").
		Joins("JOIN knowledge_subdomains AS s ON s.id = f.subdomain_id").
		Joins("JOIN knowledge_domains AS d ON d.id = s.domain_id").
		Where("f.is_active = ? AND s.is_active = ? AND d.is_active = ?", true, true, true)
	if strings.TrimSpace(subdomainID) != "" { q = q.Where("f.subdomain_id = ?", subdomainID) }
	if strings.TrimSpace(search) != "" { q = q.Where("f.name ILIKE ?", "%"+strings.TrimSpace(search)+"%") }
	err := q.Order("f.name ASC").Limit(normalizeLimit(limit)).Scan(&result).Error
	return result, err
}
