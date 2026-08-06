package projectrepository

import (
	"errors"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *model.Project) error
	GetByID(id uuid.UUID) (*model.Project, error)
	List(search, status, category string, page, limit int) ([]model.Project, int64, error)
	Update(project *model.Project) error
	Delete(id uuid.UUID) error
	CountByOwner(ownerID uuid.UUID, status string) (int64, error)
	ListByUser(userID uuid.UUID) ([]model.Project, error)

	// Applications
	CreateApplication(app *model.ProjectApplication) error
	GetApplicationByID(id uuid.UUID) (*model.ProjectApplication, error)
	ListApplications(projectID uuid.UUID) ([]model.ProjectApplication, error)
	UpdateApplication(app *model.ProjectApplication) error

	// Members
	AddMember(member *model.ProjectMember) error
	ListMembers(projectID uuid.UUID) ([]model.ProjectMember, error)
	GetMember(projectID, userID uuid.UUID) (*model.ProjectMember, error)

	// Milestones
	CreateMilestone(milestone *model.ProjectMilestone) error
	GetMilestoneByID(id uuid.UUID) (*model.ProjectMilestone, error)
	UpdateMilestone(milestone *model.ProjectMilestone) error
	ListMilestones(projectID uuid.UUID) ([]model.ProjectMilestone, error)
	GetNextMilestone(projectID uuid.UUID) (*model.ProjectMilestone, error)
}

type pgProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &pgProjectRepository{db: db}
}

func (r *pgProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *pgProjectRepository) GetByID(id uuid.UUID) (*model.Project, error) {
	var project model.Project
	err := r.db.Preload("Owner.Profile").
		Preload("Members.User.Profile").
		Preload("Milestones").
		Where("id = ?", id).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *pgProjectRepository) List(search, status, category string, page, limit int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{}).Preload("Owner.Profile")

	if status != "" {
		query = query.Where("projects.status = ?", status)
	}

	if category != "" {
		query = query.Where("projects.category = ?", category)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("projects.title ILIKE ? OR projects.description ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("projects.created_at desc").Find(&projects).Error
	return projects, total, err
}

func (r *pgProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *pgProjectRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Project{}).Error
}

func (r *pgProjectRepository) CountByOwner(ownerID uuid.UUID, status string) (int64, error) {
	var count int64
	query := r.db.Model(&model.Project{}).Where("owner_id = ?", ownerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *pgProjectRepository) ListByUser(userID uuid.UUID) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Preload("Owner.Profile").Preload("Milestones").
		Where("owner_id = ? OR id IN (SELECT project_id FROM project_members WHERE user_id = ? AND status = 'active')", userID, userID).
		Where("status IN ('open', 'active', 'in_review')").
		Order("updated_at desc").Find(&projects).Error
	return projects, err
}

// --- Applications ---

func (r *pgProjectRepository) CreateApplication(app *model.ProjectApplication) error {
	return r.db.Create(app).Error
}

func (r *pgProjectRepository) GetApplicationByID(id uuid.UUID) (*model.ProjectApplication, error) {
	var app model.ProjectApplication
	err := r.db.Preload("Applicant.Profile").Where("id = ?", id).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (r *pgProjectRepository) ListApplications(projectID uuid.UUID) ([]model.ProjectApplication, error) {
	var apps []model.ProjectApplication
	err := r.db.Preload("Applicant.Profile").
		Where("project_id = ?", projectID).
		Order("created_at desc").Find(&apps).Error
	return apps, err
}

func (r *pgProjectRepository) UpdateApplication(app *model.ProjectApplication) error {
	return r.db.Save(app).Error
}

// --- Members ---

func (r *pgProjectRepository) AddMember(member *model.ProjectMember) error {
	return r.db.Create(member).Error
}

func (r *pgProjectRepository) ListMembers(projectID uuid.UUID) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	err := r.db.Preload("User.Profile").
		Where("project_id = ?", projectID).
		Order("joined_at asc").Find(&members).Error
	return members, err
}

func (r *pgProjectRepository) GetMember(projectID, userID uuid.UUID) (*model.ProjectMember, error) {
	var member model.ProjectMember
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

// --- Milestones ---

func (r *pgProjectRepository) CreateMilestone(milestone *model.ProjectMilestone) error {
	return r.db.Create(milestone).Error
}

func (r *pgProjectRepository) GetMilestoneByID(id uuid.UUID) (*model.ProjectMilestone, error) {
	var milestone model.ProjectMilestone
	err := r.db.Where("id = ?", id).First(&milestone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &milestone, nil
}

func (r *pgProjectRepository) UpdateMilestone(milestone *model.ProjectMilestone) error {
	return r.db.Save(milestone).Error
}

func (r *pgProjectRepository) ListMilestones(projectID uuid.UUID) ([]model.ProjectMilestone, error) {
	var milestones []model.ProjectMilestone
	err := r.db.Where("project_id = ?", projectID).Order("due_at asc NULLS LAST").Find(&milestones).Error
	return milestones, err
}

func (r *pgProjectRepository) GetNextMilestone(projectID uuid.UUID) (*model.ProjectMilestone, error) {
	var milestone model.ProjectMilestone
	err := r.db.Where("project_id = ? AND status != 'completed'", projectID).
		Order("due_at asc NULLS LAST").First(&milestone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &milestone, nil
}
