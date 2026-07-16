package service

import (
	"errors"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/repository"

	"github.com/google/uuid"
)

type CollaborationService interface {
	Create(userID uuid.UUID, req dto.CreateCollaborationRequest) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*dto.CollaborationResponse, error)
	List(search, roleRequired, status string, page, limit int) ([]dto.CollaborationResponse, int64, error)
	Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateCollaborationRequest) error
	Delete(userID uuid.UUID, id uuid.UUID) error
	Apply(userID uuid.UUID, id uuid.UUID, req dto.ApplyCollaborationRequest) (uuid.UUID, error)
	ListApplications(userID uuid.UUID, id uuid.UUID) ([]dto.CollaborationApplicationResponse, error)
	UpdateApplicationStatus(userID uuid.UUID, id uuid.UUID, appID uuid.UUID, req dto.UpdateApplicationRequest) error
}

type collaborationService struct {
	collabRepo repository.CollaborationRepository
	userRepo   repository.UserRepository
}

func NewCollaborationService(collabRepo repository.CollaborationRepository, userRepo repository.UserRepository) CollaborationService {
	return &collaborationService{collabRepo: collabRepo, userRepo: userRepo}
}

func (s *collaborationService) Create(userID uuid.UUID, req dto.CreateCollaborationRequest) (uuid.UUID, error) {
	collab := &model.Collaboration{
		UserID:       userID,
		Title:        req.Title,
		Description:  req.Description,
		RoleRequired: req.RoleRequired,
		Status:       "open",
	}

	if err := s.collabRepo.Create(collab); err != nil {
		return uuid.Nil, err
	}

	return collab.ID, nil
}

func (s *collaborationService) GetByID(id uuid.UUID) (*dto.CollaborationResponse, error) {
	collab, err := s.collabRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if collab == nil {
		return nil, errors.New("collaboration proposal not found")
	}

	return &dto.CollaborationResponse{
		ID:           collab.ID,
		Title:        collab.Title,
		Description:  collab.Description,
		RoleRequired: collab.RoleRequired,
		Status:       collab.Status,
		CreatedAt:    collab.CreatedAt,
		Owner: dto.CollaborationOwner{
			ID:          collab.User.ID,
			NamaLengkap: collab.User.Profile.NamaLengkap,
			Role:        collab.User.Role,
			Afiliasi:    collab.User.Profile.Afiliasi,
			AvatarURL:   collab.User.Profile.AvatarURL,
		},
	}, nil
}

func (s *collaborationService) List(search, roleRequired, status string, page, limit int) ([]dto.CollaborationResponse, int64, error) {
	collabs, total, err := s.collabRepo.List(search, roleRequired, status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.CollaborationResponse, len(collabs))
	for i, c := range collabs {
		res[i] = dto.CollaborationResponse{
			ID:           c.ID,
			Title:        c.Title,
			Description:  c.Description,
			RoleRequired: c.RoleRequired,
			Status:       c.Status,
			CreatedAt:    c.CreatedAt,
			Owner: dto.CollaborationOwner{
				ID:          c.User.ID,
				NamaLengkap: c.User.Profile.NamaLengkap,
				Role:        c.User.Role,
				Afiliasi:    c.User.Profile.Afiliasi,
				AvatarURL:   c.User.Profile.AvatarURL,
			},
		}
	}

	return res, total, nil
}

func (s *collaborationService) Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateCollaborationRequest) error {
	collab, err := s.collabRepo.GetByID(id)
	if err != nil {
		return err
	}
	if collab == nil {
		return errors.New("collaboration proposal not found")
	}

	// Auth check: only owner can update
	if collab.UserID != userID {
		return errors.New("unauthorized to update this collaboration proposal")
	}

	collab.Title = req.Title
	collab.Description = req.Description
	collab.RoleRequired = req.RoleRequired
	collab.Status = req.Status

	return s.collabRepo.Update(collab)
}

func (s *collaborationService) Delete(userID uuid.UUID, id uuid.UUID) error {
	collab, err := s.collabRepo.GetByID(id)
	if err != nil {
		return err
	}
	if collab == nil {
		return errors.New("collaboration proposal not found")
	}

	// Auth check: only owner can delete
	if collab.UserID != userID {
		return errors.New("unauthorized to delete this collaboration proposal")
	}

	return s.collabRepo.Delete(id)
}

func (s *collaborationService) Apply(userID uuid.UUID, id uuid.UUID, req dto.ApplyCollaborationRequest) (uuid.UUID, error) {
	// Fetch collaboration proposal
	collab, err := s.collabRepo.GetByID(id)
	if err != nil {
		return uuid.Nil, err
	}
	if collab == nil {
		return uuid.Nil, errors.New("collaboration proposal not found")
	}

	// Business validation: cannot apply to own proposal
	if collab.UserID == userID {
		return uuid.Nil, errors.New("cannot apply to your own collaboration proposal")
	}

	// Business validation: project must be open
	if collab.Status != "open" {
		return uuid.Nil, errors.New("collaboration proposal is no longer open")
	}

	// Fetch applicant profile to check role match
	applicant, err := s.userRepo.GetByID(userID)
	if err != nil {
		return uuid.Nil, err
	}
	if applicant == nil {
		return uuid.Nil, errors.New("applicant not found")
	}

	// Role matching validation
	if applicant.Role != collab.RoleRequired {
		return uuid.Nil, errors.New("your role does not match the required role for this collaboration")
	}

	app := &model.Application{
		CollaborationID: id,
		UserID:          userID,
		Message:         req.Message,
		Status:          "pending",
	}

	if err := s.collabRepo.Apply(app); err != nil {
		// Can return duplicate key error if user already applied
		return uuid.Nil, errors.New("you have already applied to this collaboration proposal")
	}

	return app.ID, nil
}

func (s *collaborationService) ListApplications(userID uuid.UUID, id uuid.UUID) ([]dto.CollaborationApplicationResponse, error) {
	// Fetch collaboration to verify owner
	collab, err := s.collabRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if collab == nil {
		return nil, errors.New("collaboration proposal not found")
	}

	// Auth check: only owner can view applications
	if collab.UserID != userID {
		return nil, errors.New("unauthorized to view applications for this proposal")
	}

	apps, err := s.collabRepo.ListApplications(id)
	if err != nil {
		return nil, err
	}

	res := make([]dto.CollaborationApplicationResponse, len(apps))
	for i, a := range apps {
		res[i] = dto.CollaborationApplicationResponse{
			ID:        a.ID,
			Message:   a.Message,
			Status:    a.Status,
			CreatedAt: a.CreatedAt,
			Applicant: dto.ApplicantResponse{
				ID:          a.User.ID,
				NamaLengkap: a.User.Profile.NamaLengkap,
				Role:        a.User.Role,
				Afiliasi:    a.User.Profile.Afiliasi,
				Lokasi:      a.User.Profile.Lokasi,
				AvatarURL:   a.User.Profile.AvatarURL,
			},
		}
	}

	return res, nil
}

func (s *collaborationService) UpdateApplicationStatus(userID uuid.UUID, id uuid.UUID, appID uuid.UUID, req dto.UpdateApplicationRequest) error {
	// Verify collaboration owner
	collab, err := s.collabRepo.GetByID(id)
	if err != nil {
		return err
	}
	if collab == nil {
		return errors.New("collaboration proposal not found")
	}

	// Auth check: only owner can accept/reject applications
	if collab.UserID != userID {
		return errors.New("unauthorized to manage applications for this proposal")
	}

	app, err := s.collabRepo.GetApplicationByID(appID)
	if err != nil {
		return err
	}
	if app == nil || app.CollaborationID != id {
		return errors.New("application not found for this proposal")
	}

	app.Status = req.Status
	return s.collabRepo.UpdateApplication(app)
}
