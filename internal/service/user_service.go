package service

import (
	"errors"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserProfileResponse, error)
	GetPublicProfile(userID uuid.UUID) (*dto.UserProfileStatsResponse, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) error
	UpdateAvatar(userID uuid.UUID, avatarURL string) error
	ListUsers(search, role string, page, limit int) ([]dto.UserProfileResponse, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
	connRepo repository.ConnectionRepository
	projRepo repository.ProjectRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	connRepo repository.ConnectionRepository,
	projRepo repository.ProjectRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		connRepo: connRepo,
		projRepo: projRepo,
	}
}

func (s *userService) GetProfile(userID uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return toUserProfileResponse(user), nil
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

	profile := toUserProfileResponse(user)
	return &dto.UserProfileStatsResponse{
		UserProfileResponse: *profile,
		Stats: dto.ProfileStats{
			ConnectionCount: connCount,
			ProjectCount:    projCount,
			DiscussionCount: discCount,
		},
	}, nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) error {
	profile := &model.Profile{
		UserID:         userID,
		NamaLengkap:    req.NamaLengkap,
		GelarDepan:     req.GelarDepan,
		GelarBelakang:  req.GelarBelakang,
		Afiliasi:       req.Afiliasi,
		Lokasi:         req.Lokasi,
		BidangKeahlian: req.BidangKeahlian,
		Industry:       req.Industry,
		Bio:            req.Bio,
		Mission:        req.Mission,
		AvatarURL:      req.AvatarURL,
	}

	return s.userRepo.UpdateProfile(profile)
}

func (s *userService) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	return s.userRepo.UpdateAvatar(userID, avatarURL)
}

func (s *userService) ListUsers(search, role string, page, limit int) ([]dto.UserProfileResponse, int64, error) {
	users, total, err := s.userRepo.ListUsers(search, role, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.UserProfileResponse, len(users))
	for i, u := range users {
		res[i] = *toUserProfileResponse(&u)
	}

	return res, total, nil
}

func toUserProfileResponse(u *model.User) *dto.UserProfileResponse {
	return &dto.UserProfileResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		Profile: dto.ProfileResponse{
			NamaLengkap:    u.Profile.NamaLengkap,
			GelarDepan:     u.Profile.GelarDepan,
			GelarBelakang:  u.Profile.GelarBelakang,
			Afiliasi:       u.Profile.Afiliasi,
			Lokasi:         u.Profile.Lokasi,
			BidangKeahlian: u.Profile.BidangKeahlian,
			Industry:       u.Profile.Industry,
			Bio:            u.Profile.Bio,
			Mission:        u.Profile.Mission,
			AvatarURL:      u.Profile.AvatarURL,
			ViewCount:      u.Profile.ViewCount,
		},
	}
}
