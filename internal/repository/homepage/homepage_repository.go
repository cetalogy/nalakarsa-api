package homepagerepository

import (
	"errors"

	"nalakarsa/internal/model/homepage"

	"gorm.io/gorm"
)

type HomepageRepository interface {
	GetHero() (*homepage.HomepageHero, error)
	ListSections() ([]homepage.HomepageSection, error)
	ListTestimonials() ([]homepage.HomepageTestimonial, error)
}

type pgHomepageRepository struct {
	db *gorm.DB
}

func NewHomepageRepository(db *gorm.DB) HomepageRepository {
	return &pgHomepageRepository{db: db}
}

func (r *pgHomepageRepository) GetHero() (*homepage.HomepageHero, error) {
	var hero homepage.HomepageHero
	err := r.db.Where("is_active = ?", true).Order("updated_at desc").First(&hero).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hero, nil
}

func (r *pgHomepageRepository) ListSections() ([]homepage.HomepageSection, error) {
	var sections []homepage.HomepageSection
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (r *pgHomepageRepository) ListTestimonials() ([]homepage.HomepageTestimonial, error) {
	var testimonials []homepage.HomepageTestimonial
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&testimonials).Error
	if err != nil {
		return nil, err
	}
	return testimonials, nil
}
