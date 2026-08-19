package projectrepository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	projectcommon "nalakarsa/internal/common/project"
	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(p *model.Project) error
	GetByID(id uuid.UUID) (*model.Project, error)
	List(search, status, category string, page, limit int) ([]model.Project, int64, error)
	Update(p *model.Project) error
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

	// Collaboration Requests (FE Contract)
	CreateCollabRequest(req *model.CollaborationRequest) error
	GetCollabRequestByID(id uuid.UUID) (*model.CollaborationRequest, error)
	ListCollabRequestsForUser(userID uuid.UUID) ([]model.CollaborationRequest, error)
	HasPendingCollabRequest(applicantID uuid.UUID, discussionID, projectID *uuid.UUID) (bool, error)
	ApproveCollabRequest(requestID, initiatorID uuid.UUID) (*model.CollaborationRequest, *model.Project, *model.GroupChat, error)
	RejectCollabRequest(requestID, initiatorID uuid.UUID) (*model.CollaborationRequest, error)
}

type pgProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &pgProjectRepository{db: db}
}

func (r *pgProjectRepository) Create(p *model.Project) error {
	return r.db.Create(p).Error
}

func (r *pgProjectRepository) GetByID(id uuid.UUID) (*model.Project, error) {
	var p model.Project
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

func (r *pgProjectRepository) List(search, status, category string, page, limit int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{}).Preload("Owner")

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

func (r *pgProjectRepository) Update(p *model.Project) error {
	return r.db.Save(p).Error
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
	err := r.db.Preload("Owner").Preload("Milestones").
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
	err := r.db.Preload("Applicant").Where("id = ?", id).First(&app).Error
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
	err := r.db.Preload("Applicant").
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
	err := r.db.Preload("User").
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

func (r *pgProjectRepository) CreateCollabRequest(req *model.CollaborationRequest) error {
	return r.db.Create(req).Error
}

func (r *pgProjectRepository) GetCollabRequestByID(id uuid.UUID) (*model.CollaborationRequest, error) {
	var req model.CollaborationRequest
	err := r.db.Preload("Applicant").
		Preload("Discussion").
		Preload("Project").
		Where("id = ?", id).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *pgProjectRepository) HasPendingCollabRequest(applicantID uuid.UUID, discussionID, projectID *uuid.UUID) (bool, error) {
	var count int64
	q := r.db.Model(&model.CollaborationRequest{}).
		Where("applicant_id = ? AND status = 'PENDING'", applicantID)

	if discussionID != nil {
		q = q.Where("discussion_id = ?", *discussionID)
	}
	if projectID != nil {
		q = q.Where("project_id = ?", *projectID)
	}

	err := q.Count(&count).Error
	return count > 0, err
}

func (r *pgProjectRepository) ListCollabRequestsForUser(userID uuid.UUID) ([]model.CollaborationRequest, error) {
	var requests []model.CollaborationRequest

	// Fetch requests where the user is either the applicant OR the initiator/owner of the target discussion or project
	err := r.db.Preload("Applicant").
		Preload("Discussion").
		Preload("Project").
		Joins("LEFT JOIN discussions ON discussions.id = collaboration_requests.discussion_id").
		Joins("LEFT JOIN projects ON projects.id = collaboration_requests.project_id").
		Where("collaboration_requests.applicant_id = ? OR discussions.user_id = ? OR projects.owner_id = ?", userID, userID, userID).
		Order("collaboration_requests.created_at desc").
		Find(&requests).Error

	return requests, err
}

func (r *pgProjectRepository) ApproveCollabRequest(requestID, initiatorID uuid.UUID) (*model.CollaborationRequest, *model.Project, *model.GroupChat, error) {
	var collabReq model.CollaborationRequest
	if err := r.db.Preload("Applicant").Preload("Discussion").Preload("Project").Where("id = ?", requestID).First(&collabReq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("collaboration request not found")
		}
		return nil, nil, nil, err
	}

	if collabReq.Status != "PENDING" {
		return nil, nil, nil, errors.New("collaboration request has already been processed")
	}

	// Verify initiator authorization
	var isAuthorized bool
	if collabReq.DiscussionID != nil {
		var disc model.Discussion
		if err := r.db.Where("id = ?", *collabReq.DiscussionID).First(&disc).Error; err == nil {
			if disc.UserID == initiatorID {
				isAuthorized = true
			}
		}
	}
	if !isAuthorized && collabReq.ProjectID != nil {
		var proj model.Project
		if err := r.db.Where("id = ?", *collabReq.ProjectID).First(&proj).Error; err == nil {
			if proj.OwnerID == initiatorID {
				isAuthorized = true
			}
		}
	}

	if !isAuthorized {
		return nil, nil, nil, errors.New("only project/topic initiator can approve collaboration requests")
	}

	var targetProject model.Project
	var targetGroupChat model.GroupChat

	// ATOMIC TRANSACTION WORKFLOW
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update collaboration_requests status -> ACCEPTED
		collabReq.Status = projectcommon.CollabStatusAccepted
		if err := tx.Save(&collabReq).Error; err != nil {
			return err
		}

		// 2. Update discussions is_in_collaboration -> TRUE
		if collabReq.DiscussionID != nil {
			if err := tx.Model(&model.Discussion{}).
				Where("id = ?", *collabReq.DiscussionID).
				Update("is_in_collaboration", true).Error; err != nil {
				return err
			}
		}

		// 3. Find or Create Project with status 'Mitra Terkonfirmasi'
		var proj model.Project
		var projectFound bool

		if collabReq.ProjectID != nil {
			if err := tx.Where("id = ?", *collabReq.ProjectID).First(&proj).Error; err == nil {
				projectFound = true
			}
		} else if collabReq.DiscussionID != nil {
			if err := tx.Where("source_discussion_id = ?", *collabReq.DiscussionID).First(&proj).Error; err == nil {
				projectFound = true
			}
		}

		if projectFound {
			proj.Status = "Partner Confirmed"
			if err := tx.Save(&proj).Error; err != nil {
				return err
			}
		} else {
			// Create new project based on discussion
			var disc model.Discussion
			if collabReq.DiscussionID != nil {
				_ = tx.Where("id = ?", *collabReq.DiscussionID).First(&disc)
			}
			title := "Collaboration Project"
			desc := collabReq.ProposedContribution
			category := "Technology"
			if disc.Title != "" {
				title = disc.Title
				desc = disc.Description
				category = disc.Category
			}

			proj = model.Project{
				OwnerID:            initiatorID,
				Title:              title,
				Description:        desc,
				Category:           category,
				Needs:              "Practitioner",
				Status:             "Partner Confirmed",
				Progress:           0,
				SourceDiscussionID: collabReq.DiscussionID,
			}
			if err := tx.Create(&proj).Error; err != nil {
				return err
			}
		}

		// Link project ID to collaboration request if not set
		if collabReq.ProjectID == nil {
			collabReq.ProjectID = &proj.ID
			tx.Model(&model.CollaborationRequest{}).Where("id = ?", collabReq.ID).Update("project_id", proj.ID)
		}
		targetProject = proj

		// Add project members
		// Initiator as leader/owner
		var memberInitiator model.ProjectMember
		if err := tx.Where("project_id = ? AND user_id = ?", proj.ID, initiatorID).First(&memberInitiator).Error; err != nil {
			tx.Create(&model.ProjectMember{
				ProjectID: proj.ID,
				UserID:    initiatorID,
				Role:      "Initiator",
				Status:    "active",
			})
		}
		// Applicant as confirmed partner
		var memberApplicant model.ProjectMember
		if err := tx.Where("project_id = ? AND user_id = ?", proj.ID, collabReq.ApplicantID).First(&memberApplicant).Error; err != nil {
			tx.Create(&model.ProjectMember{
				ProjectID: proj.ID,
				UserID:    collabReq.ApplicantID,
				Role:      "Collaboration Partner",
				Status:    "active",
			})
		}

		// 4. Find or Create Group Chat for Project / Topic
		var groupChat model.GroupChat
		var groupChatFound bool

		if collabReq.DiscussionID != nil {
			if err := tx.Where("topic_id = ?", *collabReq.DiscussionID).First(&groupChat).Error; err == nil {
				groupChatFound = true
			}
		}
		if !groupChatFound && proj.ID != uuid.Nil {
			if err := tx.Where("project_id = ?", proj.ID).First(&groupChat).Error; err == nil {
				groupChatFound = true
			}
		}

		var applicantUser model.User
		applicantName := "A new collaborator"
		if err := tx.Where("id = ?", collabReq.ApplicantID).First(&applicantUser).Error; err == nil {
			if applicantUser.FullName != "" {
				applicantName = applicantUser.FullName
			} else if applicantUser.FirstName != "" {
				applicantName = strings.TrimSpace(applicantUser.FirstName + " " + applicantUser.LastName)
			}
		}

		now := time.Now()
		welcomeMsg := fmt.Sprintf("Welcome %s to the Project Collaboration Group!", applicantName)

		if !groupChatFound {
			groupChat = model.GroupChat{
				TopicID:         collabReq.DiscussionID,
				ProjectID:       &proj.ID,
				Title:           proj.Title,
				Badge:           "Collaboration Group",
				LastMessage:     &welcomeMsg,
				LastMessageTime: &now,
			}
			if err := tx.Create(&groupChat).Error; err != nil {
				return err
			}
		} else {
			groupChat.LastMessage = &welcomeMsg
			groupChat.LastMessageTime = &now
			tx.Save(&groupChat)
		}
		targetGroupChat = groupChat

		// Add Group Chat Members
		var gmInitiator model.GroupChatMember
		if err := tx.Where("group_chat_id = ? AND user_id = ?", groupChat.ID, initiatorID).First(&gmInitiator).Error; err != nil {
			tx.Create(&model.GroupChatMember{
				GroupChatID: groupChat.ID,
				UserID:      initiatorID,
				Role:        "Initiator",
			})
		}

		var gmApplicant model.GroupChatMember
		if err := tx.Where("group_chat_id = ? AND user_id = ?", groupChat.ID, collabReq.ApplicantID).First(&gmApplicant).Error; err != nil {
			tx.Create(&model.GroupChatMember{
				GroupChatID: groupChat.ID,
				UserID:      collabReq.ApplicantID,
				Role:        "Collaboration Partner",
			})
		}

		// Add initial welcome system message in group_messages
		sysMessage := model.GroupMessage{
			GroupChatID:     groupChat.ID,
			SenderID:        nil,
			IsSystemMessage: true,
			Content:         welcomeMsg,
		}
		if err := tx.Create(&sysMessage).Error; err != nil {
			return err
		}

		// 5. Send Notification to Applicant
		notif := model.Notification{
			UserID:       collabReq.ApplicantID,
			Type:         "collaboration",
			ActorID:      &initiatorID,
			ResourceType: "project",
			ResourceID:   &proj.ID,
			Payload:      `{"title":"Collaboration Request Approved","message":"Your collaboration request has been approved by the initiator."}`,
		}
		if err := tx.Create(&notif).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return &collabReq, &targetProject, &targetGroupChat, nil
}

func (r *pgProjectRepository) RejectCollabRequest(requestID, initiatorID uuid.UUID) (*model.CollaborationRequest, error) {
	var collabReq model.CollaborationRequest
	if err := r.db.Preload("Applicant").Preload("Discussion").Preload("Project").Where("id = ?", requestID).First(&collabReq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("collaboration request not found")
		}
		return nil, err
	}

	if collabReq.Status != "PENDING" {
		return nil, errors.New("collaboration request has already been processed")
	}

	// Verify initiator authorization
	var isAuthorized bool
	if collabReq.DiscussionID != nil {
		var disc model.Discussion
		if err := r.db.Where("id = ?", *collabReq.DiscussionID).First(&disc).Error; err == nil {
			if disc.UserID == initiatorID {
				isAuthorized = true
			}
		}
	}
	if !isAuthorized && collabReq.ProjectID != nil {
		var proj model.Project
		if err := r.db.Where("id = ?", *collabReq.ProjectID).First(&proj).Error; err == nil {
			if proj.OwnerID == initiatorID {
				isAuthorized = true
			}
		}
	}

	if !isAuthorized {
		return nil, errors.New("only project/topic initiator can reject collaboration requests")
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		collabReq.Status = projectcommon.CollabStatusRejected
		//collabReq.RejectionReason = &reason
		if err := tx.Save(&collabReq).Error; err != nil {
			return err
		}

		// Insert notification to applicant
		notif := model.Notification{
			UserID:       collabReq.ApplicantID,
			Type:         "collaboration",
			ActorID:      &initiatorID,
			ResourceType: "collaboration_request",
			ResourceID:   &collabReq.ID,
			Payload:      `{"title":"Collaboration Request Rejected","message":"Your collaboration request was not approved."}`,
		}
		return tx.Create(&notif).Error
	})

	if err != nil {
		return nil, err
	}

	return &collabReq, nil
}
