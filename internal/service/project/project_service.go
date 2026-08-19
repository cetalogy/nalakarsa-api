package projectservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	notificationcommon "nalakarsa/internal/common/notification"
	projectcommon "nalakarsa/internal/common/project"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	discussionrepository "nalakarsa/internal/repository/discussion"
	notificationrepository "nalakarsa/internal/repository/notification"
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
	RejectCollabRequest(requestID, initiatorID uuid.UUID) (*dto.RejectCollaborationResponse, error)
}

type projectService struct {
	projRepo  projectrepository.ProjectRepository
	userRepo  userrepository.UserRepository
	discRepo  discussionrepository.DiscussionRepository
	notifRepo notificationrepository.NotificationRepository
}

func NewProjectService(
	projRepo projectrepository.ProjectRepository,
	userRepo userrepository.UserRepository,
	discRepo discussionrepository.DiscussionRepository,
	notifRepo notificationrepository.NotificationRepository,
) ProjectService {
	return &projectService{
		projRepo:  projRepo,
		userRepo:  userRepo,
		discRepo:  discRepo,
		notifRepo: notifRepo,
	}
}

func (s *projectService) Create(userID uuid.UUID, req dto.CreateProjectRequest) (uuid.UUID, error) {
	p := &model.Project{
		OwnerID:       userID,
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Status:        projectcommon.ProjectStatusDraft,
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
		Status:             projectcommon.ProjectStatusDraft,
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

	if p.Status != projectcommon.ProjectStatusOpen {
		return uuid.Nil, errors.New("project is not open for applications")
	}

	app := &model.ProjectApplication{
		ProjectID:   projectID,
		ApplicantID: userID,
		Message:     req.Message,
		Status:      projectcommon.ApplicationStatusPending,
	}

	if err := s.projRepo.CreateApplication(app); err != nil {
		return uuid.Nil, errors.New("you have already applied to this project")
	}

	// Send in-app notification to PROJECT OWNER
	applicant, _ := s.userRepo.GetByID(userID)
	applicantName := "Someone"
	if applicant != nil && applicant.FullName != "" {
		applicantName = applicant.FullName
	}

	notifPayload, _ := json.Marshal(map[string]interface{}{
		"title":   "New Project Application",
		"message": fmt.Sprintf("%s has applied to join project '%s'", applicantName, p.Title),
	})

	notif := model.Notification{
		UserID:       p.OwnerID, // Sent to PROJECT OWNER
		Type:         notificationcommon.TypeCollaboration,
		ActorID:      &userID, // Applicant user
		ResourceType: notificationcommon.ResourceProject,
		ResourceID:   &p.ID,
		Payload:      string(notifPayload),
	}
	_ = s.notifRepo.Create(&notif)

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
	if req.Status == projectcommon.ApplicationStatusAccepted {
		// Add applicant to project members
		member := &model.ProjectMember{
			ProjectID: projectID,
			UserID:    app.ApplicantID,
			Role:      "member",
			Status:    projectcommon.MemberStatusActive,
			JoinedAt:  time.Now(),
		}
		if err := s.projRepo.AddMember(member); err != nil {
			return err
		}
	}

	// Send in-app notification to APPLICANT regarding their status update
	notifTitle := "Project Application Accepted"
	if req.Status == projectcommon.ApplicationStatusRejected {
		notifTitle = "Project Application Rejected"
	}
	notifPayload, _ := json.Marshal(map[string]interface{}{
		"title":   notifTitle,
		"message": fmt.Sprintf("Your application status for project '%s' has been updated to: %s", p.Title, req.Status),
	})

	notif := model.Notification{
		UserID:       app.ApplicantID, // Sent to APPLICANT
		Type:         notificationcommon.TypeCollaboration,
		ActorID:      &userID, // Project Owner / Initiator
		ResourceType: notificationcommon.ResourceProject,
		ResourceID:   &p.ID,
		Payload:      string(notifPayload),
	}
	_ = s.notifRepo.Create(&notif)

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
		Status:     projectcommon.MilestoneStatusPending,
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

	if req.Status == projectcommon.MilestoneStatusCompleted && milestone.CompletedAt == nil {
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
	var targetOwnerID uuid.UUID

	// Validate target resource & prevent applying to own discussion/project
	if req.DiscussionID != nil {
		disc, err := s.discRepo.GetByID(*req.DiscussionID)
		if err != nil || disc == nil {
			return nil, errors.New("discussion topic not found")
		}
		if disc.UserID == applicantID {
			return nil, errors.New("cannot apply for collaboration on your own discussion topic")
		}
		targetOwnerID = disc.UserID
		title = disc.Title
	} else if req.ProjectID != nil {
		proj, err := s.projRepo.GetByID(*req.ProjectID)
		if err != nil || proj == nil {
			return nil, errors.New("project not found")
		}
		if proj.OwnerID == applicantID {
			return nil, errors.New("cannot apply for collaboration on your own project")
		}
		targetOwnerID = proj.OwnerID
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
		Status:               projectcommon.CollabStatusPending,
	}

	if err := s.projRepo.CreateCollabRequest(&collabReq); err != nil {
		return nil, err
	}

	// Send in-app notification to OWNER / INITIATOR of the target topic or project
	notifPayload, _ := json.Marshal(map[string]interface{}{
		"title":   "New Collaboration Request",
		"message": fmt.Sprintf("%s has submitted a collaboration request on '%s'", applicant.FullName, title),
	})

	notif := model.Notification{
		UserID:       targetOwnerID, // Sent to OWNER / INITIATOR
		Type:         notificationcommon.TypeCollaboration,
		ActorID:      &applicantID, // Applicant user
		ResourceType: notificationcommon.ResourceCollaborationRequest,
		ResourceID:   &collabReq.ID,
		Payload:      string(notifPayload),
	}
	_ = s.notifRepo.Create(&notif)

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
		title := "Collaboration Project"
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
	projectTitle := "Collaboration Project"
	if proj != nil {
		projectID = &proj.ID
		if proj.Title != "" {
			projectTitle = proj.Title
		}
	}

	var groupChatID uuid.UUID
	if groupChat != nil {
		groupChatID = groupChat.ID
	}

	// Send in-app notification to the APPLICANT (user who submitted request)
	if s.notifRepo != nil {
		notifPayload, _ := json.Marshal(map[string]interface{}{
			"title":   "Collaboration Request Approved",
			"message": fmt.Sprintf("Your collaboration request on '%s' has been approved! You have been added to the Project Collaboration Group.", projectTitle),
		})

		notif := model.Notification{
			UserID:       req.ApplicantID, // Sent to APPLICANT
			Type:         notificationcommon.TypeCollaboration,
			ActorID:      &initiatorID, // Project Initiator
			ResourceType: notificationcommon.ResourceProject,
			ResourceID:   projectID,
			Payload:      string(notifPayload),
		}
		_ = s.notifRepo.Create(&notif)
	}

	return &dto.ApproveCollaborationResponse{
		RequestID:   req.ID,
		Status:      req.Status,
		ProjectID:   projectID,
		GroupChatID: groupChatID,
	}, nil
}

func (s *projectService) RejectCollabRequest(requestID, initiatorID uuid.UUID) (*dto.RejectCollaborationResponse, error) {
	collabReq, err := s.projRepo.RejectCollabRequest(requestID, initiatorID)
	if err != nil {
		notifPayload, _ := json.Marshal(map[string]interface{}{
			"title":   "Collaboration Request Rejected",
			"message": "Your collaboration request was rejected.",
		})

		notif := model.Notification{
			UserID:       collabReq.ApplicantID, // Sent to APPLICANT
			Type:         notificationcommon.TypeCollaboration,
			ActorID:      &initiatorID, // Project Initiator
			ResourceType: notificationcommon.ResourceProject,
			ResourceID:   &requestID, // Using requestID as resourceID for rejection context
			Payload:      string(notifPayload),
		}
		_ = s.notifRepo.Create(&notif)

		return &dto.RejectCollaborationResponse{
			RequestID: collabReq.ID,
			Status:    collabReq.Status,
		}, err
	}

	reasonStr := ""
	if collabReq.RejectionReason != nil {
		reasonStr = *collabReq.RejectionReason
	}

	// Send in-app notification to the APPLICANT (user who submitted request)
	if s.notifRepo != nil {
		msg := "Your collaboration request was declined."
		if reasonStr != "" {
			msg = fmt.Sprintf("Your collaboration request was declined: %s", reasonStr)
		}

		notifPayload, _ := json.Marshal(map[string]interface{}{
			"title":   "Collaboration Request Declined",
			"message": msg,
		})

		notif := model.Notification{
			UserID:       collabReq.ApplicantID, // Sent to APPLICANT
			Type:         notificationcommon.TypeCollaboration,
			ActorID:      &initiatorID, // Project Initiator
			ResourceType: notificationcommon.ResourceCollaborationRequest,
			ResourceID:   &collabReq.ID,
			Payload:      string(notifPayload),
		}
		_ = s.notifRepo.Create(&notif)
	}

	return &dto.RejectCollaborationResponse{
		RequestID: collabReq.ID,
		Status:    collabReq.Status,
		//Reason:    reasonStr,
	}, nil
}
