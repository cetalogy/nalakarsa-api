package conversationservice

import (
	"errors"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	conversationrepository "nalakarsa/internal/repository/conversation"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type ConversationService interface {
	GetOrCreateDirect(userID uuid.UUID, req dto.CreateDirectConversationRequest) (*dto.ConversationResponse, error)
	ListConversations(userID uuid.UUID, page, limit int) ([]dto.ConversationResponse, int64, error)
	ListMessages(userID uuid.UUID, conversationID uuid.UUID, limit int, cursor string) ([]dto.MessageResponse, bool, error)
	SendMessage(userID uuid.UUID, conversationID uuid.UUID, req dto.SendMessageRequest) (*dto.MessageResponse, error)
	MarkRead(userID uuid.UUID, conversationID uuid.UUID) error
}

type conversationService struct {
	convRepo conversationrepository.ConversationRepository
	userRepo userrepository.UserRepository
}

func NewConversationService(convRepo conversationrepository.ConversationRepository, userRepo userrepository.UserRepository) ConversationService {
	return &conversationService{convRepo: convRepo, userRepo: userRepo}
}

func (s *conversationService) GetOrCreateDirect(userID uuid.UUID, req dto.CreateDirectConversationRequest) (*dto.ConversationResponse, error) {
	if userID == req.TargetUserID {
		return nil, errors.New("cannot create conversation with yourself")
	}

	// Check target user exists
	target, err := s.userRepo.GetByID(req.TargetUserID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, errors.New("target user not found")
	}

	// Check if direct conversation already exists
	existing, err := s.convRepo.GetDirectByPair(userID, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return s.buildConversationResponse(existing, userID)
	}

	// Create new conversation
	conv := &model.Conversation{
		Type: "direct",
	}
	if err := s.convRepo.Create(conv); err != nil {
		return nil, err
	}

	// Add both members
	if err := s.convRepo.AddMember(&model.ConversationMember{ConversationID: conv.ID, UserID: userID}); err != nil {
		return nil, err
	}
	if err := s.convRepo.AddMember(&model.ConversationMember{ConversationID: conv.ID, UserID: req.TargetUserID}); err != nil {
		return nil, err
	}

	// Reload with preloads
	conv, err = s.convRepo.GetByID(conv.ID)
	if err != nil {
		return nil, err
	}

	return s.buildConversationResponse(conv, userID)
}

func (s *conversationService) ListConversations(userID uuid.UUID, page, limit int) ([]dto.ConversationResponse, int64, error) {
	convs, total, err := s.convRepo.ListByUser(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.ConversationResponse, len(convs))
	for i, c := range convs {
		r, err := s.buildConversationResponse(&c, userID)
		if err != nil {
			continue
		}
		res[i] = *r
	}

	return res, total, nil
}

func (s *conversationService) ListMessages(userID uuid.UUID, conversationID uuid.UUID, limit int, cursor string) ([]dto.MessageResponse, bool, error) {
	// Verify membership
	member, err := s.convRepo.GetMember(conversationID, userID)
	if err != nil {
		return nil, false, err
	}
	if member == nil {
		return nil, false, errors.New("not a member of this conversation")
	}

	if limit <= 0 {
		limit = 20
	}

	messages, err := s.convRepo.ListMessages(conversationID, limit+1, cursor)
	if err != nil {
		return nil, false, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	res := make([]dto.MessageResponse, len(messages))
	for i, m := range messages {
		res[i] = dto.MessageResponse{
			ID:        m.ID,
			Body:      m.Body,
			SenderID:  m.SenderID,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
		}
	}

	return res, hasMore, nil
}

func (s *conversationService) SendMessage(userID uuid.UUID, conversationID uuid.UUID, req dto.SendMessageRequest) (*dto.MessageResponse, error) {
	// Verify membership
	member, err := s.convRepo.GetMember(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("not a member of this conversation")
	}

	msg := &model.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Body:           req.Body,
		Status:         "sent",
	}

	if err := s.convRepo.CreateMessage(msg); err != nil {
		return nil, err
	}

	// Update conversation last_message_at
	_ = s.convRepo.UpdateLastMessageAt(conversationID)

	return &dto.MessageResponse{
		ID:        msg.ID,
		Body:      msg.Body,
		SenderID:  msg.SenderID,
		Status:    msg.Status,
		CreatedAt: msg.CreatedAt,
	}, nil
}

func (s *conversationService) MarkRead(userID uuid.UUID, conversationID uuid.UUID) error {
	member, err := s.convRepo.GetMember(conversationID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("not a member of this conversation")
	}

	// Get latest message to mark as last read
	messages, err := s.convRepo.ListMessages(conversationID, 1, "")
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		return nil
	}

	return s.convRepo.UpdateLastRead(conversationID, userID, messages[0].ID)
}

func (s *conversationService) buildConversationResponse(conv *model.Conversation, currentUserID uuid.UUID) (*dto.ConversationResponse, error) {
	// Find the other participant
	var participant dto.ConversationParticipant
	for _, m := range conv.Members {
		if m.UserID != currentUserID {
			participant = dto.ConversationParticipant{
				ID:          m.User.ID,
				NamaLengkap: m.User.Profile.NamaLengkap,
				Role:        m.User.Role,
				AvatarURL:   m.User.Profile.AvatarURL,
			}
			break
		}
	}

	// Get last message
	var lastMessage *dto.LastMessageResponse
	messages, err := s.convRepo.ListMessages(conv.ID, 1, "")
	if err == nil && len(messages) > 0 {
		lastMessage = &dto.LastMessageResponse{
			ID:        messages[0].ID,
			Body:      messages[0].Body,
			SenderID:  messages[0].SenderID,
			CreatedAt: messages[0].CreatedAt,
		}
	}

	unreadCount, _ := s.convRepo.CountUnread(conv.ID, currentUserID)

	return &dto.ConversationResponse{
		ID:          conv.ID,
		Participant: participant,
		LastMessage: lastMessage,
		UnreadCount: unreadCount,
	}, nil
}
