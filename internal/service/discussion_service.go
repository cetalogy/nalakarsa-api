package service

import (
	"errors"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/repository"

	"github.com/google/uuid"
)

type DiscussionService interface {
	Create(userID uuid.UUID, req dto.CreateDiscussionRequest) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*dto.DiscussionDetailResponse, error)
	List(search, category, role, sort string, page, limit int) ([]dto.DiscussionResponse, int64, error)
	Update(userID uuid.UUID, id uuid.UUID, req dto.UpdateDiscussionRequest) error
	Delete(userID uuid.UUID, id uuid.UUID) error
	AddComment(userID uuid.UUID, discussionID uuid.UUID, req dto.CreateCommentRequest) (*dto.DiscussionCommentResponse, error)
	DeleteComment(userID uuid.UUID, commentID uuid.UUID) error
}

type discussionService struct {
	discRepo repository.DiscussionRepository
	userRepo repository.UserRepository
}

func NewDiscussionService(discRepo repository.DiscussionRepository, userRepo repository.UserRepository) DiscussionService {
	return &discussionService{discRepo: discRepo, userRepo: userRepo}
}

func (s *discussionService) Create(userID uuid.UUID, req dto.CreateDiscussionRequest) (uuid.UUID, error) {
	disc := &model.Discussion{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
	}

	if err := s.discRepo.Create(disc); err != nil {
		return uuid.Nil, err
	}

	return disc.ID, nil
}

func (s *discussionService) GetByID(id uuid.UUID) (*dto.DiscussionDetailResponse, error) {
	disc, err := s.discRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if disc == nil {
		return nil, errors.New("discussion not found")
	}

	commentsRes := make([]dto.DiscussionCommentResponse, len(disc.Comments))
	for i, c := range disc.Comments {
		commentsRes[i] = dto.DiscussionCommentResponse{
			ID:        c.ID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			Creator: dto.DiscussionCreator{
				ID:          c.User.ID,
				NamaLengkap: c.User.Profile.NamaLengkap,
				Role:        c.User.Role,
				AvatarURL:   c.User.Profile.AvatarURL,
			},
		}
	}

	return &dto.DiscussionDetailResponse{
		ID:        disc.ID,
		Title:     disc.Title,
		Content:   disc.Content,
		Category:  disc.Category,
		Tags:      disc.Tags,
		CreatedAt: disc.CreatedAt,
		Creator: dto.DiscussionCreator{
			ID:          disc.User.ID,
			NamaLengkap: disc.User.Profile.NamaLengkap,
			Role:        disc.User.Role,
			AvatarURL:   disc.User.Profile.AvatarURL,
		},
		Comments: commentsRes,
	}, nil
}

func (s *discussionService) List(search, category, role, sort string, page, limit int) ([]dto.DiscussionResponse, int64, error) {
	discs, total, err := s.discRepo.List(search, category, role, sort, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.DiscussionResponse, len(discs))
	for i, d := range discs {
		res[i] = dto.DiscussionResponse{
			ID:        d.ID,
			Title:     d.Title,
			Content:   d.Content,
			Category:  d.Category,
			Tags:      d.Tags,
			CreatedAt: d.CreatedAt,
			Creator: dto.DiscussionCreator{
				ID:          d.User.ID,
				NamaLengkap: d.User.Profile.NamaLengkap,
				Role:        d.User.Role,
				AvatarURL:   d.User.Profile.AvatarURL,
			},
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

	// Auth check: only owner can update
	if disc.UserID != userID {
		return errors.New("unauthorized to update this discussion")
	}

	disc.Title = req.Title
	disc.Content = req.Content
	disc.Category = req.Category
	disc.Tags = req.Tags

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

	// Auth check: only owner can delete
	if disc.UserID != userID {
		return errors.New("unauthorized to delete this discussion")
	}

	return s.discRepo.Delete(id)
}

func (s *discussionService) AddComment(userID uuid.UUID, discussionID uuid.UUID, req dto.CreateCommentRequest) (*dto.DiscussionCommentResponse, error) {
	// Verify discussion exists
	disc, err := s.discRepo.GetByID(discussionID)
	if err != nil {
		return nil, err
	}
	if disc == nil {
		return nil, errors.New("discussion not found")
	}

	comment := &model.Comment{
		DiscussionID: discussionID,
		UserID:       userID,
		Content:      req.Content,
	}

	if err := s.discRepo.CreateComment(comment); err != nil {
		return nil, err
	}

	// Fetch user details for response
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.DiscussionCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		Creator: dto.DiscussionCreator{
			ID:          user.ID,
			NamaLengkap: user.Profile.NamaLengkap,
			Role:        user.Role,
			AvatarURL:   user.Profile.AvatarURL,
		},
	}, nil
}

func (s *discussionService) DeleteComment(userID uuid.UUID, commentID uuid.UUID) error {
	comment, err := s.discRepo.GetCommentByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Verify discussion to see if current user is owner of discussion or comment
	disc, err := s.discRepo.GetByID(comment.DiscussionID)
	if err != nil {
		return err
	}

	// Auth check: owner of comment OR owner of discussion can delete comments
	if comment.UserID != userID && (disc == nil || disc.UserID != userID) {
		return errors.New("unauthorized to delete this comment")
	}

	return s.discRepo.DeleteComment(commentID)
}
