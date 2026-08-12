package projectservice

import (
	"errors"
	"time"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	discussionrepository "nalakarsa/internal/repository/discussion"
	projectrepository "nalakarsa/internal/repository/project"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(userID uuid.UUID, req dto.CreateProjectRequest) (uuid.UUID, error)
	CreateFromDiscussion(userID uuid.UUID, discussionID uuid.UUID) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*dto.ProjectDetailResponse, error)
	List(search, status, category string, page, limit int) ([]dto.ProjectResponse, int64, error)
	Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateProjectRequest) error
	Delete(userID uuid.UUID, id uuid.UUID) error
	Apply(userID uuid.UUID, projectID uuid.UUID, req dto.ApplyProjectRequest) (uuid.UUID, error)
	ListApplications(userID uuid.UUID, projectID uuid.UUID) ([]dto.ProjectApplicationResponse, error)
	UpdateApplicationStatus(userID uuid.UUID, projectID uuid.UUID, appID uuid.UUID, req dto.UpdateApplicationStatusRequest) error
	ListMembers(projectID uuid.UUID) ([]dto.ProjectMemberResponse, error)
	CreateMilestone(userID uuid.UUID, projectID uuid.UUID, req dto.CreateMilestoneRequest) (uuid.UUID, error)
	UpdateMilestone(userID uuid.UUID, projectID uuid.UUID, milestoneID uuid.UUID, req dto.UpdateMilestoneRequest) error
}

type projectService struct {
	projRepo projectrepository.ProjectRepository
	userRepo userrepository.UserRepository
	discRepo discussionrepository.DiscussionRepository
}

func NewProjectService(projRepo projectrepository.ProjectRepository, userRepo userrepository.UserRepository, discRepo discussionrepository.DiscussionRepository) ProjectService {
	return &projectService{projRepo: projRepo, userRepo: userRepo, discRepo: discRepo}
}

func (s *projectService) Create(userID uuid.UUID, req dto.CreateProjectRequest) (uuid.UUID, error) {
	project := &model.Project{
		OwnerID:       userID,
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Status:        "draft",
		Needs:  req.Needs,
		FundingStatus: req.FundingStatus,
		Location:      req.Location,
		Deadline:      req.Deadline,
		Progress:      0,
	}

	if err := s.projRepo.Create(project); err != nil {
		return uuid.Nil, err
	}

	return project.ID, nil
}

func (s *projectService) CreateFromDiscussion(userID uuid.UUID, discussionID uuid.UUID) (uuid.UUID, error) {
	disc, err := s.discRepo.GetByID(discussionID)
	if err != nil {
		return uuid.Nil, err
	}
	if disc == nil {
		return uuid.Nil, errors.New("discussion not found")
	}

	project := &model.Project{
		OwnerID:       userID,
		Title:         disc.Title,
		Description:   disc.Description,
		Category:      disc.Category,
		Status:        "draft",
		Needs:       "Praktisi", // Default or map it somehow
		FundingStatus:      "Belum Ada",
		Location:           "Remote",
		Progress:           0,
		SourceDiscussionID: &discussionID,
	}

	if err := s.projRepo.Create(project); err != nil {
		return uuid.Nil, err
	}

	return project.ID, nil
}

func (s *projectService) GetByID(id uuid.UUID) (*dto.ProjectDetailResponse, error) {
	project, err := s.projRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	members := make([]dto.ProjectMemberResponse, len(project.Members))
	for i, m := range project.Members {
		members[i] = dto.ProjectMemberResponse{
			ID:          m.ID,
			UserID:      m.UserID,
			FullName:    m.User.FullName,
			Role:        m.Role,
			Status:      m.Status,
			JoinedAt:    m.JoinedAt,
			AvatarURL:   m.User.AvatarURL,
		}
	}

	milestones := make([]dto.ProjectMilestoneResponse, len(project.Milestones))
	for i, ms := range project.Milestones {
		milestones[i] = dto.ProjectMilestoneResponse{
			ID:          ms.ID,
			Title:       ms.Title,
			DueAt:       ms.DueAt,
			Status:      ms.Status,
			AssigneeID:  ms.AssigneeID,
			CompletedAt: ms.CompletedAt,
			CreatedAt:   ms.CreatedAt,
		}
	}

	return &dto.ProjectDetailResponse{
		ProjectResponse: toProjectResponse(project),
		Members:         members,
		Milestones:      milestones,
	}, nil
}

func (s *projectService) List(search, status, category string, page, limit int) ([]dto.ProjectResponse, int64, error) {
	projects, total, err := s.projRepo.List(search, status, category, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		res[i] = toProjectResponse(&p)
	}

	return res, total, nil
}

func (s *projectService) Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateProjectRequest) error {
	project, err := s.projRepo.GetByID(id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	if project.OwnerID != userID {
		return errors.New("unauthorized to update this project")
	}

	project.Title = req.Title
	project.Description = req.Description
	project.Category = req.Category
	project.Status = req.Status
	project.Needs = req.Needs
	project.FundingStatus = req.FundingStatus
	project.Location = req.Location
	project.Deadline = req.Deadline
	if req.Progress != nil {
		project.Progress = *req.Progress
	}

	return s.projRepo.Update(project)
}

func (s *projectService) Delete(userID uuid.UUID, id uuid.UUID) error {
	project, err := s.projRepo.GetByID(id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	if project.OwnerID != userID {
		return errors.New("unauthorized to delete this project")
	}

	return s.projRepo.Delete(id)
}

func (s *projectService) Apply(userID uuid.UUID, projectID uuid.UUID, req dto.ApplyProjectRequest) (uuid.UUID, error) {
	project, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return uuid.Nil, err
	}
	if project == nil {
		return uuid.Nil, errors.New("project not found")
	}

	if project.OwnerID == userID {
		return uuid.Nil, errors.New("cannot apply to your own project")
	}

	if project.Status != "open" {
		return uuid.Nil, errors.New("project is not open for applications")
	}

	app := &model.ProjectApplication{
		ProjectID:   projectID,
		ApplicantID: userID,
		Message:     req.Message,
		Status:      "pending",
	}

	if err := s.projRepo.CreateApplication(app); err != nil {
		return uuid.Nil, errors.New("you have already applied to this project")
	}

	return app.ID, nil
}

func (s *projectService) ListApplications(userID uuid.UUID, projectID uuid.UUID) ([]dto.ProjectApplicationResponse, error) {
	project, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	if project.OwnerID != userID {
		return nil, errors.New("unauthorized to view applications for this project")
	}

	apps, err := s.projRepo.ListApplications(projectID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ProjectApplicationResponse, len(apps))
	for i, a := range apps {
		res[i] = dto.ProjectApplicationResponse{
			ID:        a.ID,
			Message:   a.Message,
			Status:    a.Status,
			CreatedAt: a.CreatedAt,
			Applicant: dto.ApplicantResponse{
				ID:          a.Applicant.ID,
				FullName:    a.Applicant.FullName,
				Role:        a.Applicant.Role,
				Afiliasi:    a.Applicant.Affiliation,
				Lokasi:      a.Applicant.Location,
				AvatarURL:   a.Applicant.AvatarURL,
			},
		}
	}

	return res, nil
}

func (s *projectService) UpdateApplicationStatus(userID uuid.UUID, projectID uuid.UUID, appID uuid.UUID, req dto.UpdateApplicationStatusRequest) error {
	project, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	if project.OwnerID != userID {
		return errors.New("unauthorized to manage applications for this project")
	}

	app, err := s.projRepo.GetApplicationByID(appID)
	if err != nil {
		return err
	}
	if app == nil || app.ProjectID != projectID {
		return errors.New("application not found for this project")
	}

	app.Status = req.Status

	if err := s.projRepo.UpdateApplication(app); err != nil {
		return err
	}

	// If accepted, add as project member
	if req.Status == "accepted" {
		member := &model.ProjectMember{
			ProjectID: projectID,
			UserID:    app.ApplicantID,
			Role:      "member",
			Status:    "active",
			JoinedAt:  time.Now(),
		}
		_ = s.projRepo.AddMember(member)
	}

	return nil
}

func (s *projectService) ListMembers(projectID uuid.UUID) ([]dto.ProjectMemberResponse, error) {
	members, err := s.projRepo.ListMembers(projectID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ProjectMemberResponse, len(members))
	for i, m := range members {
		res[i] = dto.ProjectMemberResponse{
			ID:          m.ID,
			UserID:      m.UserID,
			FullName:    m.User.FullName,
			Role:        m.Role,
			Status:      m.Status,
			JoinedAt:    m.JoinedAt,
			AvatarURL:   m.User.AvatarURL,
		}
	}

	return res, nil
}

func (s *projectService) CreateMilestone(userID uuid.UUID, projectID uuid.UUID, req dto.CreateMilestoneRequest) (uuid.UUID, error) {
	project, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return uuid.Nil, err
	}
	if project == nil {
		return uuid.Nil, errors.New("project not found")
	}

	if project.OwnerID != userID {
		return uuid.Nil, errors.New("unauthorized to create milestones for this project")
	}

	milestone := &model.ProjectMilestone{
		ProjectID:  projectID,
		Title:      req.Title,
		DueAt:      req.DueAt,
		Status:     "pending",
		AssigneeID: req.AssigneeID,
	}

	if err := s.projRepo.CreateMilestone(milestone); err != nil {
		return uuid.Nil, err
	}

	return milestone.ID, nil
}

func (s *projectService) UpdateMilestone(userID uuid.UUID, projectID uuid.UUID, milestoneID uuid.UUID, req dto.UpdateMilestoneRequest) error {
	project, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	if project.OwnerID != userID {
		return errors.New("unauthorized to update milestones for this project")
	}

	milestone, err := s.projRepo.GetMilestoneByID(milestoneID)
	if err != nil {
		return err
	}
	if milestone == nil || milestone.ProjectID != projectID {
		return errors.New("milestone not found for this project")
	}

	milestone.Title = req.Title
	milestone.DueAt = req.DueAt
	milestone.Status = req.Status
	milestone.AssigneeID = req.AssigneeID

	if req.Status == "completed" && milestone.CompletedAt == nil {
		now := time.Now()
		milestone.CompletedAt = &now
	}

	return s.projRepo.UpdateMilestone(milestone)
}

func toProjectResponse(p *model.Project) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:                 p.ID,
		Title:              p.Title,
		Description:        p.Description,
		Category:           p.Category,
		Status:             p.Status,
		Needs:              p.Needs,
		FundingStatus:      p.FundingStatus,
		Location:           p.Location,
		Deadline:           p.Deadline,
		Progress:           p.Progress,
		CreatedAt:          p.CreatedAt,
		Initiator:          p.Owner.FullName,
		SourceDiscussionID: p.SourceDiscussionID,
	}
}
