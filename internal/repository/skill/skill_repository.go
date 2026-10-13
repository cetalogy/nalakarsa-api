package skill

import (
	"strings"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Search(search string, limit int) ([]model.Skill, error)
	Create(skill *model.Skill) (*model.Skill, error)
	GetByID(id uuid.UUID) (*model.Skill, error)
	ListByUser(userID uuid.UUID) ([]model.Skill, error)
	AddToUser(userID, skillID uuid.UUID) error
	RemoveFromUser(userID, skillID uuid.UUID) error
}

type repository struct{ db *gorm.DB }

func NewSkillRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Search(search string, limit int) ([]model.Skill, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	query := r.db.Where("is_active = ?", true)
	if value := strings.TrimSpace(search); value != "" {
		query = query.Where("name ILIKE ?", "%"+value+"%")
	}
	var items []model.Skill
	err := query.Order("name ASC").Limit(limit).Find(&items).Error
	return items, err
}

func (r *repository) Create(item *model.Skill) (*model.Skill, error) {
	var existing model.Skill
	if err := r.db.Where("lower(name) = lower(?)", item.Name).First(&existing).Error; err == nil {
		return &existing, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
		return nil, err
	}
	return r.findByName(item.Name)
}

func (r *repository) findByName(name string) (*model.Skill, error) {
	var item model.Skill
	err := r.db.Where("lower(name) = lower(?)", name).First(&item).Error
	return &item, err
}

func (r *repository) GetByID(id uuid.UUID) (*model.Skill, error) {
	var item model.Skill
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&item).Error
	return &item, err
}

func (r *repository) ListByUser(userID uuid.UUID) ([]model.Skill, error) {
	var items []model.Skill
	err := r.db.Table("skills").Joins("JOIN user_skills ON user_skills.skill_id = skills.id").
		Where("user_skills.user_id = ? AND skills.is_active = ?", userID, true).
		Order("skills.name ASC").Find(&items).Error
	return items, err
}

func (r *repository) AddToUser(userID, skillID uuid.UUID) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.UserSkill{
		UserID: userID, SkillID: skillID,
	}).Error
}

func (r *repository) RemoveFromUser(userID, skillID uuid.UUID) error {
	result := r.db.Where("user_id = ? AND skill_id = ?", userID, skillID).Delete(&model.UserSkill{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
