package discussionservice

import (
	"errors"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	discussionrepository "nalakarsa/internal/repository/discussion"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type DiscussionService interface {
	Create(userID uuid.UUID, req dto.CreateDiscussionRequest) (uuid.UUID, error)
	GetByID(id uuid.UUID, currentUserID *uuid.UUID) (*dto.DiscussionDetailResponse, error)
	List(search, category, role, sort string, page, limit int, currentUserID *uuid.UUID) ([]dto.DiscussionResponse, int64, error)
	Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateDiscussionRequest) error
	Delete(userID uuid.UUID, id uuid.UUID) error
	AddReply(userID uuid.UUID, discussionID uuid.UUID, req dto.CreateReplyRequest) (*dto.DiscussionReplyResponse, error)
	ListReplies(discussionID uuid.UUID, page, limit int) (*dto.DiscussionRepliesData, int64, error)
	DeleteReply(userID uuid.UUID, replyID uuid.UUID) error
	Vote(userID uuid.UUID, discussionID uuid.UUID) error
	Unvote(userID uuid.UUID, discussionID uuid.UUID) error
	MarkCollaboration(userID uuid.UUID, id uuid.UUID) error
}

type discussionService struct {
	discRepo discussionrepository.DiscussionRepository
	userRepo userrepository.UserRepository
}

func NewDiscussionService(discRepo discussionrepository.DiscussionRepository, userRepo userrepository.UserRepository) DiscussionService {
	return &discussionService{discRepo: discRepo, userRepo: userRepo}
}

func (s *discussionService) Create(userID uuid.UUID, req dto.CreateDiscussionRequest) (uuid.UUID, error) {
	disc := &model.Discussion{
		UserID:      userID,
		Title:       req.Title,
		Description: req.GetDescription(),
		Category:    req.Category,
		Tags:        req.Tags,
		Status:      "open",
	}

	if err := s.discRepo.Create(disc); err != nil {
		return uuid.Nil, err
	}

	return disc.ID, nil
}

func (s *discussionService) GetByID(id uuid.UUID, currentUserID *uuid.UUID) (*dto.DiscussionDetailResponse, error) {
	disc, err := s.discRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if disc == nil {
		return nil, errors.New("discussion not found")
	}

	replyCount, _ := s.discRepo.CountReplies(id)
	upvoteCount, _ := s.discRepo.CountVotes(id)

	hasUpvoted := false
	if currentUserID != nil {
		hasUpvoted, _ = s.discRepo.HasVoted(*currentUserID, id)
	}

	return &dto.DiscussionDetailResponse{
		ID:                 disc.ID,
		Title:              disc.Title,
		Description:        disc.Description,
		Excerpt:            disc.Description,
		Category:           disc.Category,
		Tags:               disc.Tags,
		Status:             disc.Status,
		IsInCollaboration:  disc.IsInCollaboration,
		Replies:            replyCount,
		UpvoteCount:        upvoteCount,
		HasUpvoted:         hasUpvoted,
		Time:               disc.CreatedAt,
		CreatedAt:          disc.CreatedAt,
		Author:             disc.User.FullName,
		Role:               disc.User.Role,
		SourceDiscussionID: disc.SourceDiscussionID,
		Creator: dto.DiscussionCreator{
			ID:        disc.User.ID,
			FullName:  disc.User.FullName,
			Role:      disc.User.Role,
			AvatarURL: disc.User.AvatarURL,
		},
		RepliesList: s.toDiscussionReplyResponses(disc.Replies),
	}, nil
}

func (s *discussionService) ListReplies(discussionID uuid.UUID, page, limit int) (*dto.DiscussionRepliesData, int64, error) {
	disc, err := s.discRepo.GetByID(discussionID)
	if err != nil {
		return nil, 0, err
	}
	if disc == nil {
		return nil, 0, errors.New("discussion not found")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	replies, total, err := s.discRepo.ListReplies(discussionID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	replyResponses := s.toDiscussionReplyResponses(replies)
	if replyResponses == nil {
		replyResponses = []dto.DiscussionReplyResponse{}
	}

	return &dto.DiscussionRepliesData{
		DiscussionID: disc.ID,
		TopicTitle:   disc.Title,
		TotalReplies: total,
		Replies:      replyResponses,
	}, total, nil
}

func (s *discussionService) List(search, category, role, sort string, page, limit int, currentUserID *uuid.UUID) ([]dto.DiscussionResponse, int64, error) {
	discs, total, err := s.discRepo.List(search, category, role, sort, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.DiscussionResponse, len(discs))
	for i, d := range discs {
		replyCount, _ := s.discRepo.CountReplies(d.ID)
		upvoteCount, _ := s.discRepo.CountVotes(d.ID)

		hasUpvoted := false
		if currentUserID != nil {
			hasUpvoted, _ = s.discRepo.HasVoted(*currentUserID, d.ID)
		}

		res[i] = dto.DiscussionResponse{
			ID:                 d.ID,
			Title:              d.Title,
			Description:        d.Description,
			Category:           d.Category,
			Tags:               d.Tags,
			Status:             d.Status,
			IsInCollaboration:  d.IsInCollaboration,
			Replies:            replyCount,
			UpvoteCount:        upvoteCount,
			HasUpvoted:         hasUpvoted,
			Time:               d.CreatedAt,
			Author:             d.User.FullName,
			Role:               d.User.Role,
			SourceDiscussionID: d.SourceDiscussionID,
		}
	}

	return res, total, nil
}

func (s *discussionService) Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateDiscussionRequest) error {
	disc, err := s.discRepo.GetByID(id)
	if err != nil {
		return err
	}
	if disc == nil {
		return errors.New("discussion not found")
	}
	if disc.UserID != userID {
		return errors.New("unauthorized to update this discussion")
	}

	disc.Title = req.Title
	if desc := req.GetDescription(); desc != "" {
		disc.Description = desc
	}
	disc.Category = req.Category
	disc.Tags = req.Tags
	if req.Status != "" {
		disc.Status = req.Status
	}

	return s.discRepo.Update(disc)
}

func (s *discussionService) Delete(userID uuid.UUID, id uuid.UUID) error {
	disc, err := s.discRepo.GetByID(id)
	if err != nil {
		return err
	}
	if disc == nil {
		return errors.New("discussion not found")
	}
	if disc.UserID != userID {
		return errors.New("unauthorized to delete this discussion")
	}

	return s.discRepo.Delete(id)
}

func (s *discussionService) AddReply(userID uuid.UUID, discussionID uuid.UUID, req dto.CreateReplyRequest) (*dto.DiscussionReplyResponse, error) {
	disc, err := s.discRepo.GetByID(discussionID)
	if err != nil {
		return nil, err
	}
	if disc == nil {
		return nil, errors.New("discussion not found")
	}

	reply := &model.DiscussionReply{
		DiscussionID: discussionID,
		UserID:       userID,
		ParentID:     req.ParentID,
		Content:      req.Content,
	}

	if err := s.discRepo.CreateReply(reply); err != nil {
		return nil, err
	}
	u, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.DiscussionReplyResponse{
		ID:        reply.ID,
		Content:   reply.Content,
		ParentID:  reply.ParentID,
		CreatedAt: reply.CreatedAt,
		Creator: dto.DiscussionCreator{
			ID:        u.ID,
			FullName:  u.FullName,
			Role:      u.Role,
			AvatarURL: u.AvatarURL,
		},
	}, nil
}

func (s *discussionService) DeleteReply(userID uuid.UUID, replyID uuid.UUID) error {
	reply, err := s.discRepo.GetReplyByID(replyID)
	if err != nil {
		return err
	}
	if reply == nil {
		return errors.New("reply not found")
	}
	disc, err := s.discRepo.GetByID(reply.DiscussionID)
	if err != nil {
		return err
	}
	if reply.UserID != userID && (disc == nil || disc.UserID != userID) {
		return errors.New("unauthorized to delete this reply")
	}

	return s.discRepo.DeleteReply(replyID)
}

func (s *discussionService) toDiscussionReplyResponse(r model.DiscussionReply) dto.DiscussionReplyResponse {
	return dto.DiscussionReplyResponse{
		ID:        r.ID,
		Content:   r.Content,
		ParentID:  r.ParentID,
		CreatedAt: r.CreatedAt,
		Creator: dto.DiscussionCreator{
			ID:        r.User.ID,
			FullName:  r.User.FullName,
			Role:      r.User.Role,
			AvatarURL: r.User.AvatarURL,
		},
	}
}

func (s *discussionService) toDiscussionReplyResponses(replies []model.DiscussionReply) []dto.DiscussionReplyResponse {
	if len(replies) == 0 {
		return []dto.DiscussionReplyResponse{}
	}
	res := make([]dto.DiscussionReplyResponse, len(replies))
	for i, r := range replies {
		res[i] = s.toDiscussionReplyResponse(r)
	}
	return res
}

func (s *discussionService) Vote(userID uuid.UUID, discussionID uuid.UUID) error {
	disc, err := s.discRepo.GetByID(discussionID)
	if err != nil {
		return err
	}
	if disc == nil {
		return errors.New("discussion not found")
	}
	hasVoted, err := s.discRepo.HasVoted(userID, discussionID)
	if err != nil {
		return err
	}
	if hasVoted {
		return errors.New("already upvoted this discussion")
	}

	vote := &model.DiscussionVote{
		UserID:       userID,
		DiscussionID: discussionID,
	}

	return s.discRepo.CreateVote(vote)
}

func (s *discussionService) Unvote(userID uuid.UUID, discussionID uuid.UUID) error {
	return s.discRepo.DeleteVote(userID, discussionID)
}

func (s *discussionService) MarkCollaboration(userID uuid.UUID, id uuid.UUID) error {
	disc, err := s.discRepo.GetByID(id)
	if err != nil {
		return err
	}
	if disc == nil {
		return errors.New("discussion not found")
	}

	disc.IsInCollaboration = true
	return s.discRepo.Update(disc)
}
