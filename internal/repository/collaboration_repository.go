package repository

import (
	"errors"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollaborationRepository interface {
	Create(collab *model.Collaboration) error
	GetByID(id uuid.UUID) (*model.Collaboration, error)
	List(search, roleRequired, status string, page, limit int) ([]model.Collaboration, int64, error)
	Update(collab *model.Collaboration) error
	Delete(id uuid.UUID) error
	Apply(app *model.Application) error
	ListApplications(collaborationID uuid.UUID) ([]model.Application, error)
	GetApplicationByID(id uuid.UUID) (*model.Application, error)
	UpdateApplication(app *model.Application) error
}

type pgCollaborationRepository struct {
	db *gorm.DB
}

func NewCollaborationRepository(db *gorm.DB) CollaborationRepository {
	return &pgCollaborationRepository{db: db}
}

func (r *pgCollaborationRepository) Create(collab *model.Collaboration) error {
	return r.db.Create(collab).Error
}

func (r *pgCollaborationRepository) GetByID(id uuid.UUID) (*model.Collaboration, error) {
	var collab model.Collaboration
	err := r.db.Preload("User.Profile").Where("id = ?", id).First(&collab).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &collab, nil
}

func (r *pgCollaborationRepository) List(search, roleRequired, status string, page, limit int) ([]model.Collaboration, int64, error) {
	var collabs []model.Collaboration
	var total int64

	query := r.db.Model(&model.Collaboration{}).Preload("User.Profile")

	if roleRequired != "" {
		query = query.Where("role_required = ?", roleRequired)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&collabs).Error
	if err != nil {
		return nil, 0, err
	}

	return collabs, total, nil
}

func (r *pgCollaborationRepository) Update(collab *model.Collaboration) error {
	return r.db.Save(collab).Error
}

func (r *pgCollaborationRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Collaboration{}).Error
}

func (r *pgCollaborationRepository) Apply(app *model.Application) error {
	return r.db.Create(app).Error
}

func (r *pgCollaborationRepository) ListApplications(collaborationID uuid.UUID) ([]model.Application, error) {
	var apps []model.Application
	err := r.db.Preload("User.Profile").
		Where("collaboration_id = ?", collaborationID).
		Order("created_at desc").Find(&apps).Error
	return apps, err
}

func (r *pgCollaborationRepository) GetApplicationByID(id uuid.UUID) (*model.Application, error) {
	var app model.Application
	err := r.db.Where("id = ?", id).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (r *pgCollaborationRepository) UpdateApplication(app *model.Application) error {
	return r.db.Save(app).Error
}
