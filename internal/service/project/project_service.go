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

	// Collaboration Requests (FE Contract)
	SubmitCollabRequest(applicantID uuid.UUID, req dto.SubmitCollaborationRequest) (*dto.CollaborationRequestItemResponse, error)
	ListCollabRequests(userID uuid.UUID) ([]dto.CollaborationRequestItemResponse, error)
	ApproveCollabRequest(requestID, initiatorID uuid.UUID) (*dto.ApproveCollaborationResponse, error)
	RejectCollabRequest(requestID, initiatorID uuid.UUID, req dto.RejectCollaborationRequest) (*dto.RejectCollaborationResponse, error)
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
	p := &model.Project{
		OwnerID:       userID,
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Status:        "draft",
		Needs:         req.Needs,
		FundingStatus: req.FundingStatus,
		Location:      req.Location,
		Deadline:      req.Deadline,
		Progress:      0,
	}

	if err := s.projRepo.Create(p); err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (s *projectService) CreateFromDiscussion(userID uuid.UUID, discussionID uuid.UUID) (uuid.UUID, error) {
	disc, err := s.discRepo.GetByID(discussionID)
	if err != nil {
		return uuid.Nil, err
	}
	if disc == nil {
		return uuid.Nil, errors.New("discussion not found")
	}

	p := &model.Project{
		OwnerID:            userID,
		Title:              disc.Title,
		Description:        disc.Description,
		Category:           disc.Category,
		Status:             "draft",
		Needs:              "Praktisi", // Default or map it somehow
		FundingStatus:      "Belum Ada",
		Location:           "Remote",
		Progress:           0,
		SourceDiscussionID: &discussionID,
	}

	if err := s.projRepo.Create(p); err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (s *projectService) GetByID(id uuid.UUID) (*dto.ProjectDetailResponse, error) {
	p, err := s.projRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("project not found")
	}

	members := make([]dto.ProjectMemberResponse, len(p.Members))
	for i, m := range p.Members {
		members[i] = dto.ProjectMemberResponse{
			ID:        m.ID,
			UserID:    m.UserID,
			FullName:  m.User.FullName,
			Role:      m.Role,
			Status:    m.Status,
			JoinedAt:  m.JoinedAt,
			AvatarURL: m.User.AvatarURL,
		}
	}

	milestones := make([]dto.ProjectMilestoneResponse, len(p.Milestones))
	for i, ms := range p.Milestones {
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
		ProjectResponse: toProjectResponse(p),
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
	p, err := s.projRepo.GetByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("project not found")
	}

	if p.OwnerID != userID {
		return errors.New("unauthorized to update this project")
	}

	p.Title = req.Title
	p.Description = req.Description
	p.Category = req.Category
	p.Status = req.Status
	p.Needs = req.Needs
	p.FundingStatus = req.FundingStatus
	p.Location = req.Location
	p.Deadline = req.Deadline
	if req.Progress != nil {
		p.Progress = *req.Progress
	}

	return s.projRepo.Update(p)
}

func (s *projectService) Delete(userID uuid.UUID, id uuid.UUID) error {
	p, err := s.projRepo.GetByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("project not found")
	}

	if p.OwnerID != userID {
		return errors.New("unauthorized to delete this project")
	}

	return s.projRepo.Delete(id)
}

func (s *projectService) Apply(userID uuid.UUID, projectID uuid.UUID, req dto.ApplyProjectRequest) (uuid.UUID, error) {
	p, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return uuid.Nil, err
	}
	if p == nil {
		return uuid.Nil, errors.New("project not found")
	}

	if p.OwnerID == userID {
		return uuid.Nil, errors.New("cannot apply to your own project")
	}

	if p.Status != "open" {
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
	p, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("project not found")
	}

	if p.OwnerID != userID {
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
				ID:        a.Applicant.ID,
				FullName:  a.Applicant.FullName,
				Role:      a.Applicant.Role,
				Afiliasi:  a.Applicant.Affiliation,
				Lokasi:    a.Applicant.Location,
				AvatarURL: a.Applicant.AvatarURL,
			},
		}
	}

	return res, nil
}

func (s *projectService) UpdateApplicationStatus(userID uuid.UUID, projectID uuid.UUID, appID uuid.UUID, req dto.UpdateApplicationStatusRequest) error {
	p, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("project not found")
	}

	if p.OwnerID != userID {
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
			ID:        m.ID,
			UserID:    m.UserID,
			FullName:  m.User.FullName,
			Role:      m.Role,
			Status:    m.Status,
			JoinedAt:  m.JoinedAt,
			AvatarURL: m.User.AvatarURL,
		}
	}

	return res, nil
}

func (s *projectService) CreateMilestone(userID uuid.UUID, projectID uuid.UUID, req dto.CreateMilestoneRequest) (uuid.UUID, error) {
	p, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return uuid.Nil, err
	}
	if p == nil {
		return uuid.Nil, errors.New("project not found")
	}

	if p.OwnerID != userID {
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
	p, err := s.projRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("project not found")
	}

	if p.OwnerID != userID {
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

func (s *projectService) SubmitCollabRequest(applicantID uuid.UUID, req dto.SubmitCollaborationRequest) (*dto.CollaborationRequestItemResponse, error) {
	if req.DiscussionID == nil && req.ProjectID == nil {
		return nil, errors.New("discussionId or projectId is required")
	}

	applicant, err := s.userRepo.GetByID(applicantID)
	if err != nil || applicant == nil {
		return nil, errors.New("applicant user not found")
	}

	var title string

	// Validate target resource & prevent applying to own discussion/project
	if req.DiscussionID != nil {
		disc, err := s.discRepo.GetByID(*req.DiscussionID)
		if err != nil || disc == nil {
			return nil, errors.New("discussion topic not found")
		}
		if disc.UserID == applicantID {
			return nil, errors.New("cannot apply for collaboration on your own discussion topic")
		}
		title = disc.Title
	} else if req.ProjectID != nil {
		proj, err := s.projRepo.GetByID(*req.ProjectID)
		if err != nil || proj == nil {
			return nil, errors.New("project not found")
		}
		if proj.OwnerID == applicantID {
			return nil, errors.New("cannot apply for collaboration on your own project")
		}
		title = proj.Title
	}

	// Check if already has a pending application
	hasPending, err := s.projRepo.HasPendingCollabRequest(applicantID, req.DiscussionID, req.ProjectID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, errors.New("you already have a pending collaboration request for this topic/project")
	}

	collabReq := model.CollaborationRequest{
		DiscussionID:         req.DiscussionID,
		ProjectID:            req.ProjectID,
		ApplicantID:          applicantID,
		ProposedContribution: req.ProposedContribution,
		Status:               "PENDING",
	}

	if err := s.projRepo.CreateCollabRequest(&collabReq); err != nil {
		return nil, err
	}

	return &dto.CollaborationRequestItemResponse{
		ID:                   collabReq.ID,
		DiscussionID:         collabReq.DiscussionID,
		ProjectID:            collabReq.ProjectID,
		Title:                title,
		ProposedContribution: collabReq.ProposedContribution,
		Status:               collabReq.Status,
		CreatedAt:            collabReq.CreatedAt,
		Applicant: dto.CollabApplicantResponse{
			ID:        applicant.ID,
			FullName:  applicant.FullName,
			Role:      applicant.Role,
			AvatarURL: applicant.AvatarURL,
		},
	}, nil
}

func (s *projectService) ListCollabRequests(userID uuid.UUID) ([]dto.CollaborationRequestItemResponse, error) {
	requests, err := s.projRepo.ListCollabRequestsForUser(userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.CollaborationRequestItemResponse, len(requests))
	for i, r := range requests {
		title := "Proyek Kolaborasi"
		if r.Discussion != nil && r.Discussion.Title != "" {
			title = r.Discussion.Title
		} else if r.Project != nil && r.Project.Title != "" {
			title = r.Project.Title
		}

		res[i] = dto.CollaborationRequestItemResponse{
			ID:                   r.ID,
			DiscussionID:         r.DiscussionID,
			ProjectID:            r.ProjectID,
			Title:                title,
			ProposedContribution: r.ProposedContribution,
			Status:               r.Status,
			RejectionReason:      r.RejectionReason,
			CreatedAt:            r.CreatedAt,
			Applicant: dto.CollabApplicantResponse{
				ID:        r.Applicant.ID,
				FullName:  r.Applicant.FullName,
				Role:      r.Applicant.Role,
				AvatarURL: r.Applicant.AvatarURL,
			},
		}
	}

	return res, nil
}

func (s *projectService) ApproveCollabRequest(requestID, initiatorID uuid.UUID) (*dto.ApproveCollaborationResponse, error) {
	req, proj, groupChat, err := s.projRepo.ApproveCollabRequest(requestID, initiatorID)
	if err != nil {
		return nil, err
	}

	var projectID *uuid.UUID
	if proj != nil {
		projectID = &proj.ID
	}

	var groupChatID uuid.UUID
	if groupChat != nil {
		groupChatID = groupChat.ID
	}

	return &dto.ApproveCollaborationResponse{
		RequestID:   req.ID,
		Status:      req.Status,
		ProjectID:   projectID,
		GroupChatID: groupChatID,
	}, nil
}

func (s *projectService) RejectCollabRequest(requestID, initiatorID uuid.UUID, req dto.RejectCollaborationRequest) (*dto.RejectCollaborationResponse, error) {
	collabReq, err := s.projRepo.RejectCollabRequest(requestID, initiatorID, req.Reason)
	if err != nil {
		return nil, err
	}

	reasonStr := ""
	if collabReq.RejectionReason != nil {
		reasonStr = *collabReq.RejectionReason
	}

	return &dto.RejectCollaborationResponse{
		RequestID: collabReq.ID,
		Status:    collabReq.Status,
		Reason:    reasonStr,
	}, nil
}

