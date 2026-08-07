package userservice

import (
	"errors"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	connectionrepository "nalakarsa/internal/repository/connection"
	projectrepository "nalakarsa/internal/repository/project"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	GetPublicProfile(userID uuid.UUID) (*dto.UserProfileStatsResponse, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) error
	UpdateAvatar(userID uuid.UUID, avatarURL string) error
	ListUsers(search, role string, page, limit int) ([]dto.UserResponse, int64, error)
	GetMyProjects(userID uuid.UUID) (*dto.MyProjectsResponse, error)
	GetMyStats(userID uuid.UUID) (*dto.UserStatsResponse, error)
}

type userService struct {
	userRepo userrepository.UserRepository
	connRepo connectionrepository.ConnectionRepository
	projRepo projectrepository.ProjectRepository
}

func NewUserService(
	userRepo userrepository.UserRepository,
	connRepo connectionrepository.ConnectionRepository,
	projRepo projectrepository.ProjectRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		connRepo: connRepo,
		projRepo: projRepo,
	}
}

func (s *userService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return toUserResponse(user), nil
}

func (s *userService) GetPublicProfile(userID uuid.UUID) (*dto.UserProfileStatsResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Increment view count
	_ = s.userRepo.IncrementViewCount(userID)

	// Get stats
	connCount, _ := s.connRepo.CountAccepted(userID)
	projCount, _ := s.projRepo.CountByOwner(userID, "")
	discCount, _ := s.userRepo.CountDiscussions(userID)

	profile := toUserResponse(user)
	return &dto.UserProfileStatsResponse{
		UserResponse: *profile,
		Stats: dto.ProfileStats{
			ConnectionCount: connCount,
			ProjectCount:    projCount,
			DiscussionCount: discCount,
		},
	}, nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) error {
	fullName := req.FullName
	if fullName == "" {
		fullName = req.FirstName + " " + req.LastName
	}
	
	user := &model.User{
		ID:          userID,
		FirstName:   req.FirstName,
		MiddleName:  &req.MiddleName,
		LastName:    req.LastName,
		FullName:    fullName,
		PrefixTitle: req.PrefixTitle,
		SuffixTitle: req.SuffixTitle,
		Affiliation: req.Affiliation,
		Location:    req.Location,
		Expertise:   req.Expertise,
		Industry:    req.Industry,
		Bio:         req.Bio,
		Mission:     req.Mission,
		AvatarURL:   req.AvatarURL,
	}

	return s.userRepo.UpdateProfile(user)
}

func (s *userService) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	return s.userRepo.UpdateAvatar(userID, avatarURL)
}

func (s *userService) ListUsers(search, role string, page, limit int) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.ListUsers(search, role, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.UserResponse, len(users))
	for i, u := range users {
		res[i] = *toUserResponse(&u)
	}

	return res, total, nil
}

func (s *userService) GetMyProjects(userID uuid.UUID) (*dto.MyProjectsResponse, error) {
	projects, err := s.projRepo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		project := toProjectResponse(p)
		if project.Initiator == "" {
			owner, err := s.userRepo.GetByID(p.OwnerID)
			if err == nil && owner != nil {
				project.Initiator = owner.FullName
			}
		}
		res[i] = project
	}

	return &dto.MyProjectsResponse{
		Projects: res,
		Total:    int64(len(res)),
	}, nil
}

func (s *userService) GetMyStats(userID uuid.UUID) (*dto.UserStatsResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	connCount, _ := s.connRepo.CountAccepted(userID)
	projCount, _ := s.projRepo.CountByOwner(userID, "")
	discCount, _ := s.userRepo.CountDiscussions(userID)

	return &dto.UserStatsResponse{
		ConnectionCount: connCount,
		ProjectCount:    projCount,
		DiscussionCount: discCount,
		ViewCount:       int64(user.ViewCount),
	}, nil
}

func toUserResponse(u *model.User) *dto.UserResponse {
	middleName := ""
	if u.MiddleName != nil {
		middleName = *u.MiddleName
	}
	return &dto.UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		FirstName:   u.FirstName,
		MiddleName:  &middleName,
		LastName:    u.LastName,
		FullName:    u.FullName,
		PrefixTitle: u.PrefixTitle,
		SuffixTitle: u.SuffixTitle,
		Affiliation: u.Affiliation,
		Location:    u.Location,
		Expertise:   u.Expertise,
		Industry:    u.Industry,
		Mission:     u.Mission,
		AvatarURL:   u.AvatarURL,
	}
}

func toProjectResponse(p model.Project) dto.ProjectResponse {
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
