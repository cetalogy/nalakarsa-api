package homepagerepository

import (
	"nalakarsa/internal/model"
	"errors"

	"gorm.io/gorm"
)

type HomepageRepository interface {
	GetHero() (*model.HomepageHero, error)
	ListSections() ([]model.HomepageSection, error)
	ListTestimonials() ([]model.HomepageTestimonial, error)
}

type pgHomepageRepository struct {
	db *gorm.DB
}

func NewHomepageRepository(db *gorm.DB) HomepageRepository {
	return &pgHomepageRepository{db: db}
}

func (r *pgHomepageRepository) GetHero() (*model.HomepageHero, error) {
	var hero model.HomepageHero
	err := r.db.Where("is_active = ?", true).Order("updated_at desc").First(&hero).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hero, nil
}

func (r *pgHomepageRepository) ListSections() ([]model.HomepageSection, error) {
	var sections []model.HomepageSection
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (r *pgHomepageRepository) ListTestimonials() ([]model.HomepageTestimonial, error) {
	var testimonials []model.HomepageTestimonial
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&testimonials).Error
	if err != nil {
		return nil, err
	}
	return testimonials, nil
}
