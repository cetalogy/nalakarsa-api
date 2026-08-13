package conversationservice

import (
	"errors"
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model/conversation"
	"nalakarsa/internal/model/user"
	conversationrepository "nalakarsa/internal/repository/conversation"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type ConversationService interface {
	GetOrCreateDirect(userID uuid.UUID, req dto.CreateDirectConversationRequest) (*dto.ConversationResponse, error)
	StartChat(userID uuid.UUID, req dto.StartChatRequest) (*dto.ConversationResponse, error)
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
	conv := &conversation.Conversation{
		Type: "direct",
	}
	if err := s.convRepo.Create(conv); err != nil {
		return nil, err
	}

	// Add both members
	if err := s.convRepo.AddMember(&conversation.ConversationMember{ConversationID: conv.ID, UserID: userID}); err != nil {
		return nil, err
	}
	if err := s.convRepo.AddMember(&conversation.ConversationMember{ConversationID: conv.ID, UserID: req.TargetUserID}); err != nil {
		return nil, err
	}

	// Reload with preloads
	conv, err = s.convRepo.GetByID(conv.ID)
	if err != nil {
		return nil, err
	}

	return s.buildConversationResponse(conv, userID)
}

func (s *conversationService) StartChat(userID uuid.UUID, req dto.StartChatRequest) (*dto.ConversationResponse, error) {
	targetUserID, err := s.resolveTargetUserID(req.Name, req.Role)
	if err != nil {
		return nil, err
	}

	// Resolve to explicit direct-conversation target and reuse existing method.
	return s.GetOrCreateDirect(userID, dto.CreateDirectConversationRequest{
		TargetUserID: targetUserID,
	})
}

// resolveTargetUserID finds a unique active recipient by exact full name and optional role.
// If zero or multiple matches are found, it returns an explicit error for deterministic behavior.
func (s *conversationService) resolveTargetUserID(name, role string) (uuid.UUID, error) {
	cleanName := strings.TrimSpace(name)
	cleanRole := strings.TrimSpace(role)

	users, _, err := s.userRepo.ListUsers(cleanName, cleanRole, 1, 10)
	if err != nil {
		return uuid.Nil, err
	}
	if len(users) == 0 {
		return uuid.Nil, errors.New("target user not found")
	}

	candidates := make([]user.User, 0, len(users))
	for _, u := range users {
		if strings.EqualFold(strings.TrimSpace(u.FullName), cleanName) {
			candidates = append(candidates, u)
		}
	}

	switch len(candidates) {
	case 0:
		return uuid.Nil, errors.New("target user not found by exact name")
	case 1:
		return candidates[0].ID, nil
	default:
		return uuid.Nil, errors.New("multiple users found for the provided name and role, use conversations/direct with user id")
	}
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
		sender := "them"
		if m.SenderID == userID {
			sender = "me"
		}
		res[i] = dto.MessageResponse{
			ID:     m.ID,
			Sender: sender,
			Text:   m.Body,
			Time:   m.CreatedAt,
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

	text := req.Text
	if text == "" {
		text = req.Body // fallback
	}
	msg := &conversation.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Body:           text,
		Status:         "sent",
	}

	if err := s.convRepo.CreateMessage(msg); err != nil {
		return nil, err
	}

	// Update conversation last_message_at
	_ = s.convRepo.UpdateLastMessageAt(conversationID)

	sender := "them"
	if msg.SenderID == userID {
		sender = "me"
	}
	return &dto.MessageResponse{
		ID:     msg.ID,
		Sender: sender,
		Text:   msg.Body,
		Time:   msg.CreatedAt,
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

func (s *conversationService) buildConversationResponse(conv *conversation.Conversation, currentUserID uuid.UUID) (*dto.ConversationResponse, error) {
	// Find the other participant
	var name, role, avatar string
	for _, m := range conv.Members {
		if m.UserID != currentUserID {
			name = m.User.FullName
			role = m.User.Role
			avatar = m.User.AvatarURL
			break
		}
	}

	// Get last message
	var lastMessageText string
	messages, err := s.convRepo.ListMessages(conv.ID, 1, "")
	if err == nil && len(messages) > 0 {
		lastMessageText = messages[0].Body
	}

	return &dto.ConversationResponse{
		ID:          conv.ID,
		Name:        name,
		Role:        role,
		Avatar:      avatar,
		LastMessage: lastMessageText,
	}, nil
}
