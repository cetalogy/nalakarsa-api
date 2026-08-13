package projectrepository

import (
	"errors"

	"nalakarsa/internal/model/project"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(p *project.Project) error
	GetByID(id uuid.UUID) (*project.Project, error)
	List(search, status, category string, page, limit int) ([]project.Project, int64, error)
	Update(p *project.Project) error
	Delete(id uuid.UUID) error
	CountByOwner(ownerID uuid.UUID, status string) (int64, error)
	ListByUser(userID uuid.UUID) ([]project.Project, error)

	// Applications
	CreateApplication(app *project.ProjectApplication) error
	GetApplicationByID(id uuid.UUID) (*project.ProjectApplication, error)
	ListApplications(projectID uuid.UUID) ([]project.ProjectApplication, error)
	UpdateApplication(app *project.ProjectApplication) error

	// Members
	AddMember(member *project.ProjectMember) error
	ListMembers(projectID uuid.UUID) ([]project.ProjectMember, error)
	GetMember(projectID, userID uuid.UUID) (*project.ProjectMember, error)

	// Milestones
	CreateMilestone(milestone *project.ProjectMilestone) error
	GetMilestoneByID(id uuid.UUID) (*project.ProjectMilestone, error)
	UpdateMilestone(milestone *project.ProjectMilestone) error
	ListMilestones(projectID uuid.UUID) ([]project.ProjectMilestone, error)
	GetNextMilestone(projectID uuid.UUID) (*project.ProjectMilestone, error)
}

type pgProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &pgProjectRepository{db: db}
}

func (r *pgProjectRepository) Create(p *project.Project) error {
	return r.db.Create(p).Error
}

func (r *pgProjectRepository) GetByID(id uuid.UUID) (*project.Project, error) {
	var p project.Project
	err := r.db.Preload("Owner").
		Preload("Members.User").
		Preload("Milestones").
		Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgProjectRepository) List(search, status, category string, page, limit int) ([]project.Project, int64, error) {
	var projects []project.Project
	var total int64

	query := r.db.Model(&project.Project{}).Preload("Owner")

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

func (r *pgProjectRepository) Update(p *project.Project) error {
	return r.db.Save(p).Error
}

func (r *pgProjectRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&project.Project{}).Error
}

func (r *pgProjectRepository) CountByOwner(ownerID uuid.UUID, status string) (int64, error) {
	var count int64
	query := r.db.Model(&project.Project{}).Where("owner_id = ?", ownerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *pgProjectRepository) ListByUser(userID uuid.UUID) ([]project.Project, error) {
	var projects []project.Project
	err := r.db.Preload("Owner").Preload("Milestones").
		Where("owner_id = ? OR id IN (SELECT project_id FROM project_members WHERE user_id = ? AND status = 'active')", userID, userID).
		Where("status IN ('open', 'active', 'in_review')").
		Order("updated_at desc").Find(&projects).Error
	return projects, err
}

// --- Applications ---

func (r *pgProjectRepository) CreateApplication(app *project.ProjectApplication) error {
	return r.db.Create(app).Error
}

func (r *pgProjectRepository) GetApplicationByID(id uuid.UUID) (*project.ProjectApplication, error) {
	var app project.ProjectApplication
	err := r.db.Preload("Applicant").Where("id = ?", id).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (r *pgProjectRepository) ListApplications(projectID uuid.UUID) ([]project.ProjectApplication, error) {
	var apps []project.ProjectApplication
	err := r.db.Preload("Applicant").
		Where("project_id = ?", projectID).
		Order("created_at desc").Find(&apps).Error
	return apps, err
}

func (r *pgProjectRepository) UpdateApplication(app *project.ProjectApplication) error {
	return r.db.Save(app).Error
}

// --- Members ---

func (r *pgProjectRepository) AddMember(member *project.ProjectMember) error {
	return r.db.Create(member).Error
}

func (r *pgProjectRepository) ListMembers(projectID uuid.UUID) ([]project.ProjectMember, error) {
	var members []project.ProjectMember
	err := r.db.Preload("User").
		Where("project_id = ?", projectID).
		Order("joined_at asc").Find(&members).Error
	return members, err
}

func (r *pgProjectRepository) GetMember(projectID, userID uuid.UUID) (*project.ProjectMember, error) {
	var member project.ProjectMember
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

func (r *pgProjectRepository) CreateMilestone(milestone *project.ProjectMilestone) error {
	return r.db.Create(milestone).Error
}

func (r *pgProjectRepository) GetMilestoneByID(id uuid.UUID) (*project.ProjectMilestone, error) {
	var milestone project.ProjectMilestone
	err := r.db.Where("id = ?", id).First(&milestone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &milestone, nil
}

func (r *pgProjectRepository) UpdateMilestone(milestone *project.ProjectMilestone) error {
	return r.db.Save(milestone).Error
}

func (r *pgProjectRepository) ListMilestones(projectID uuid.UUID) ([]project.ProjectMilestone, error) {
	var milestones []project.ProjectMilestone
	err := r.db.Where("project_id = ?", projectID).Order("due_at asc NULLS LAST").Find(&milestones).Error
	return milestones, err
}

func (r *pgProjectRepository) GetNextMilestone(projectID uuid.UUID) (*project.ProjectMilestone, error) {
	var milestone project.ProjectMilestone
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
