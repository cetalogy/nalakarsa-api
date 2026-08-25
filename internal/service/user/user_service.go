package userservice

import (
	"encoding/json"
	"errors"
	"fmt"

	notificationcommon "nalakarsa/internal/common/notification"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	connectionrepository "nalakarsa/internal/repository/connection"
	notificationrepository "nalakarsa/internal/repository/notification"
	projectrepository "nalakarsa/internal/repository/project"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	GetPublicProfile(userID uuid.UUID) (*dto.UserProfileStatsResponse, error)
	ResolveUser(identifier string) (*model.User, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) error
	UpdateAvatar(userID uuid.UUID, avatarURL string) error
	ListUsers(search, role string, page, limit int) ([]dto.UserResponse, int64, error)
	GetMyProjects(userID uuid.UUID) (*dto.MyProjectsResponse, error)
	GetMyStats(userID uuid.UUID) (*dto.UserStatsResponse, error)
	ToggleFollow(currentUserID, targetUserID uuid.UUID) (*dto.ToggleFollowResponse, error)
	GetFollowers(currentUserID *uuid.UUID, targetUserID uuid.UUID, page, limit int) ([]dto.FollowUserItemResponse, int64, error)
	GetFollowing(currentUserID *uuid.UUID, targetUserID uuid.UUID, page, limit int) ([]dto.FollowUserItemResponse, int64, error)
}

type userService struct {
	userRepo  userrepository.UserRepository
	connRepo  connectionrepository.ConnectionRepository
	projRepo  projectrepository.ProjectRepository
	notifRepo notificationrepository.NotificationRepository
}

func NewUserService(
	userRepo userrepository.UserRepository,
	connRepo connectionrepository.ConnectionRepository,
	projRepo projectrepository.ProjectRepository,
	notifRepo notificationrepository.NotificationRepository,
) UserService {
	return &userService{
		userRepo:  userRepo,
		connRepo:  connRepo,
		projRepo:  projRepo,
		notifRepo: notifRepo,
	}
}

func (s *userService) ResolveUser(identifier string) (*model.User, error) {
	return s.userRepo.GetByIDOrIdentifier(identifier)
}

func (s *userService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	u, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	return toUserResponse(u), nil
}

func (s *userService) GetPublicProfile(userID uuid.UUID) (*dto.UserProfileStatsResponse, error) {
	u, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Increment view count
	_ = s.userRepo.IncrementViewCount(userID)

	// Get stats
	connCount, _ := s.connRepo.CountAccepted(userID)
	projCount, _ := s.projRepo.CountByOwner(userID, "")
	discCount, _ := s.userRepo.CountDiscussions(userID)

	profile := toUserResponse(u)
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
	existingUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	firstName := req.FirstName
	if firstName == "" {
		firstName = existingUser.FirstName
	}

	middleName := req.MiddleName
	var middleNamePtr *string
	if middleName != "" {
		middleNamePtr = &middleName
	} else {
		middleNamePtr = existingUser.MiddleName
	}

	lastName := req.LastName
	if lastName == "" {
		lastName = existingUser.LastName
	}

	fullName := req.FullName
	if fullName == "" {
		if middleNamePtr != nil && *middleNamePtr != "" {
			fullName = fmt.Sprintf("%s %s %s", firstName, *middleNamePtr, lastName)
		} else {
			fullName = fmt.Sprintf("%s %s", firstName, lastName)
		}
	}

	prefixTitle := req.PrefixTitle
	if prefixTitle == "" {
		prefixTitle = existingUser.PrefixTitle
	}

	suffixTitle := req.SuffixTitle
	if suffixTitle == "" {
		suffixTitle = existingUser.SuffixTitle
	}

	affiliation := req.Affiliation
	if affiliation == "" {
		affiliation = existingUser.Affiliation
	}

	location := req.Location
	if location == "" {
		location = existingUser.Location
	}

	expertise := req.Expertise
	if expertise == "" {
		expertise = existingUser.Expertise
	}

	industry := req.Industry
	if industry == "" {
		industry = existingUser.Industry
	}

	bio := req.Bio
	if bio == "" {
		bio = existingUser.Bio
	}

	mission := req.Mission
	if mission == "" {
		mission = existingUser.Mission
	}

	avatarURL := req.AvatarURL
	if avatarURL == "" {
		avatarURL = existingUser.AvatarURL
	}

	u := &model.User{
		ID:          userID,
		FirstName:   firstName,
		MiddleName:  middleNamePtr,
		LastName:    lastName,
		FullName:    fullName,
		PrefixTitle: prefixTitle,
		SuffixTitle: suffixTitle,
		Affiliation: affiliation,
		Location:    location,
		Expertise:   expertise,
		Industry:    industry,
		Bio:         bio,
		Mission:     mission,
		AvatarURL:   avatarURL,
	}

	return s.userRepo.UpdateProfile(u)
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
		projRes := toProjectResponse(p)
		if projRes.Initiator == "" {
			owner, err := s.userRepo.GetByID(p.OwnerID)
			if err == nil && owner != nil {
				projRes.Initiator = owner.FullName
			}
		}
		res[i] = projRes
	}

	return &dto.MyProjectsResponse{
		Projects: res,
		Total:    int64(len(res)),
	}, nil
}

func (s *userService) GetMyStats(userID uuid.UUID) (*dto.UserStatsResponse, error) {
	u, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	connCount, _ := s.connRepo.CountAccepted(userID)
	followersCount, _ := s.userRepo.CountFollowers(userID)
	followingCount, _ := s.userRepo.CountFollowing(userID)
	projCount, _ := s.projRepo.CountByOwner(userID, "")
	discCount, _ := s.userRepo.CountDiscussions(userID)

	return &dto.UserStatsResponse{
		ConnectionCount: connCount,
		FollowersCount:  followersCount,
		FollowingCount:  followingCount,
		ProjectCount:    projCount,
		DiscussionCount: discCount,
		ViewCount:       int64(u.ViewCount),
	}, nil
}

func (s *userService) ToggleFollow(currentUserID, targetUserID uuid.UUID) (*dto.ToggleFollowResponse, error) {
	if currentUserID == targetUserID {
		return nil, errors.New("cannot follow yourself")
	}

	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return nil, err
	}
	if targetUser == nil {
		return nil, errors.New("target user not found")
	}

	isFollowing, err := s.userRepo.ToggleFollow(currentUserID, targetUserID)
	if err != nil {
		return nil, err
	}

	message := "User followed successfully"
	if !isFollowing {
		message = "User unfollowed successfully"
	} else if s.notifRepo != nil {
		// Send in-app notification to the followed user
		follower, _ := s.userRepo.GetByID(currentUserID)
		followerName := "Someone"
		if follower != nil && follower.FullName != "" {
			followerName = follower.FullName
		}

		notifPayload, _ := json.Marshal(map[string]interface{}{
			"title":   "New Follower",
			"message": fmt.Sprintf("%s is now following you", followerName),
		})

		notif := model.Notification{
			UserID:       targetUserID, // Sent to the user being followed
			Type:         notificationcommon.TypeFollow,
			ActorID:      &currentUserID, // Follower
			ResourceType: notificationcommon.ResourceUser,
			ResourceID:   &currentUserID,
			Payload:      string(notifPayload),
		}
		_ = s.notifRepo.Create(&notif)
	}

	return &dto.ToggleFollowResponse{
		Message:      message,
		IsFollowing:  isFollowing,
		TargetUserID: targetUserID,
	}, nil
}

func (s *userService) GetFollowers(currentUserID *uuid.UUID, targetUserID uuid.UUID, page, limit int) ([]dto.FollowUserItemResponse, int64, error) {
	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return nil, 0, err
	}
	if targetUser == nil {
		return nil, 0, errors.New("user not found")
	}

	users, total, err := s.userRepo.GetFollowers(targetUserID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.FollowUserItemResponse, len(users))
	for i, u := range users {
		isFollowing := false
		if currentUserID != nil && *currentUserID != u.ID {
			isFollowing, _ = s.userRepo.IsFollowing(*currentUserID, u.ID)
		}
		res[i] = dto.FollowUserItemResponse{
			ID:          u.ID,
			FullName:    u.FullName,
			Role:        u.Role,
			Affiliation: u.Affiliation,
			AvatarURL:   u.AvatarURL,
			IsFollowing: isFollowing,
		}
	}

	return res, total, nil
}

func (s *userService) GetFollowing(currentUserID *uuid.UUID, targetUserID uuid.UUID, page, limit int) ([]dto.FollowUserItemResponse, int64, error) {
	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return nil, 0, err
	}
	if targetUser == nil {
		return nil, 0, errors.New("user not found")
	}

	users, total, err := s.userRepo.GetFollowing(targetUserID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.FollowUserItemResponse, len(users))
	for i, u := range users {
		isFollowing := false
		if currentUserID != nil && *currentUserID != u.ID {
			isFollowing, _ = s.userRepo.IsFollowing(*currentUserID, u.ID)
		}
		res[i] = dto.FollowUserItemResponse{
			ID:          u.ID,
			FullName:    u.FullName,
			Role:        u.Role,
			Affiliation: u.Affiliation,
			AvatarURL:   u.AvatarURL,
			IsFollowing: isFollowing,
		}
	}

	return res, total, nil
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
