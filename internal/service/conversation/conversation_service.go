package conversationservice

import (
	"errors"
	"fmt"
	"strings"
	"time"

	conversationcommon "nalakarsa/internal/common/conversation"
	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	conversationrepository "nalakarsa/internal/repository/conversation"
	userrepository "nalakarsa/internal/repository/user"
	"nalakarsa/internal/utils"

	"github.com/google/uuid"
)

type ConversationService interface {
	GetOrCreateDirect(userID uuid.UUID, req dto.CreateDirectConversationRequest) (*dto.ConversationResponse, error)
	StartChat(userID uuid.UUID, req dto.StartChatRequest) (*dto.ConversationResponse, error)
	ListConversations(userID uuid.UUID, page, limit int) ([]dto.ConversationResponse, int64, error)
	ListMessages(userID uuid.UUID, conversationID uuid.UUID, limit int, cursor string) ([]dto.MessageResponse, bool, error)
	SendMessage(userID uuid.UUID, conversationID uuid.UUID, req dto.SendMessageRequest) (*dto.MessageResponse, error)
	MarkRead(userID uuid.UUID, conversationID uuid.UUID) error

	// Group Chats (FE Contract Specification)
	ListGroupChats(userID uuid.UUID) ([]dto.GroupChatResponse, error)
	ListGroupMessages(userID uuid.UUID, groupChatID uuid.UUID, page, limit int) ([]dto.GroupMessageResponse, int64, error)
	SendGroupMessage(userID uuid.UUID, groupChatID uuid.UUID, req dto.SendGroupMessageRequest) (*dto.GroupMessageResponse, error)

	// Delete message (Direct & Group)
	DeleteMessage(userID uuid.UUID, messageID uuid.UUID) error
}

type conversationService struct {
	convRepo conversationrepository.ConversationRepository
	userRepo userrepository.UserRepository
	cfg      *config.Config
}

func NewConversationService(convRepo conversationrepository.ConversationRepository, userRepo userrepository.UserRepository, cfg *config.Config) ConversationService {
	return &conversationService{convRepo: convRepo, userRepo: userRepo, cfg: cfg}
}

func (s *conversationService) GetOrCreateDirect(userID uuid.UUID, req dto.CreateDirectConversationRequest) (*dto.ConversationResponse, error) {
	cleanID := strings.TrimPrefix(strings.TrimSpace(req.TargetUserID), "user_")
	targetUUID, err := uuid.Parse(cleanID)
	if err != nil {
		return nil, errors.New("invalid target user id format")
	}

	if userID == targetUUID {
		return nil, errors.New("cannot create conversation with yourself")
	}

	// Check target user exists
	target, err := s.userRepo.GetByID(targetUUID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, errors.New("target user not found")
	}

	// Check if direct conversation already exists
	existing, err := s.convRepo.GetDirectByPair(userID, targetUUID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return s.buildConversationResponse(existing, userID)
	}

	// Create new conversation
	conv := &model.Conversation{
		Type: conversationcommon.ConversationTypeDirect,
	}
	if err := s.convRepo.Create(conv); err != nil {
		return nil, err
	}

	// Add both members
	if err := s.convRepo.AddMember(&model.ConversationMember{ConversationID: conv.ID, UserID: userID}); err != nil {
		return nil, err
	}
	if err := s.convRepo.AddMember(&model.ConversationMember{ConversationID: conv.ID, UserID: targetUUID}); err != nil {
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
		TargetUserID: targetUserID.String(),
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

	candidates := make([]model.User, 0, len(users))
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
	msg := &model.Message{
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

	response := &dto.MessageResponse{
		ID:     msg.ID,
		Sender: sender,
		Text:   msg.Body,
		Time:   msg.CreatedAt,
	}

	// Trigger Firebase Realtime Database push asynchronously
	firebasePath := fmt.Sprintf("chats/direct/%s/messages", conversationID.String())
	utils.PushToFirebaseRealtime(firebasePath, response, s.cfg)

	return response, nil
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

func (s *conversationService) ListGroupChats(userID uuid.UUID) ([]dto.GroupChatResponse, error) {
	groupChats, err := s.convRepo.ListGroupChatsByUser(userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.GroupChatResponse, len(groupChats))
	for i, gc := range groupChats {
		var members []dto.GroupChatMemberResponse
		for _, m := range gc.Members {
			name := m.User.FullName
			if name == "" {
				name = m.User.FirstName + " " + m.User.LastName
			}
			members = append(members, dto.GroupChatMemberResponse{
				ID:        m.User.ID,
				Name:      strings.TrimSpace(name),
				Role:      m.Role,
				AvatarURL: m.User.AvatarURL,
			})
		}

		res[i] = dto.GroupChatResponse{
			ID:              gc.ID,
			TopicID:         gc.TopicID,
			ProjectID:       gc.ProjectID,
			Title:           gc.Title,
			Badge:           gc.Badge,
			LastMessage:     gc.LastMessage,
			LastMessageTime: gc.LastMessageTime,
			CreatedAt:       gc.CreatedAt,
			Members:         members,
		}
	}

	return res, nil
}

func (s *conversationService) ListGroupMessages(userID uuid.UUID, groupChatID uuid.UUID, page, limit int) ([]dto.GroupMessageResponse, int64, error) {
	gc, err := s.convRepo.GetGroupChatByID(groupChatID)
	if err != nil {
		return nil, 0, err
	}
	if gc == nil {
		return nil, 0, errors.New("group chat not found")
	}

	isMember, err := s.convRepo.IsGroupChatMember(gc.ID, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isMember {
		return nil, 0, errors.New("unauthorized to view messages in this group chat")
	}

	messages, total, err := s.convRepo.ListGroupMessages(gc.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.GroupMessageResponse, len(messages))
	for i, m := range messages {
		var senderName, senderRole, senderAvatar string
		if m.Sender != nil {
			senderName = m.Sender.FullName
			senderRole = m.Sender.Role
			senderAvatar = m.Sender.AvatarURL
		} else if m.IsSystemMessage {
			senderName = "Nalakarsa"
			senderRole = "System"
		}

		res[i] = dto.GroupMessageResponse{
			ID:              m.ID,
			GroupChatID:     m.GroupChatID,
			SenderID:        m.SenderID,
			SenderName:      senderName,
			SenderRole:      senderRole,
			SenderAvatar:    senderAvatar,
			IsSystemMessage: m.IsSystemMessage,
			Content:         m.Content,
			CreatedAt:       m.CreatedAt,
		}
	}

	return res, total, nil
}

func (s *conversationService) SendGroupMessage(userID uuid.UUID, groupChatID uuid.UUID, req dto.SendGroupMessageRequest) (*dto.GroupMessageResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		content = strings.TrimSpace(req.Text)
	}
	if content == "" {
		return nil, errors.New("message content is required")
	}

	gc, err := s.convRepo.GetGroupChatByID(groupChatID)
	if err != nil {
		return nil, err
	}
	if gc == nil {
		return nil, errors.New("group chat not found")
	}

	isMember, err := s.convRepo.IsGroupChatMember(gc.ID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("unauthorized to send messages to this group chat")
	}

	sender, err := s.userRepo.GetByID(userID)
	if err != nil || sender == nil {
		return nil, errors.New("sender user not found")
	}

	msg := model.GroupMessage{
		GroupChatID:     gc.ID,
		SenderID:        &userID,
		IsSystemMessage: false,
		Content:         content,
	}

	if err := s.convRepo.CreateGroupMessage(&msg); err != nil {
		return nil, err
	}

	_ = s.convRepo.UpdateGroupChatLastMessage(gc.ID, content)

	response := &dto.GroupMessageResponse{
		ID:              msg.ID,
		GroupChatID:     msg.GroupChatID,
		SenderID:        msg.SenderID,
		SenderName:      sender.FullName,
		SenderRole:      sender.Role,
		SenderAvatar:    sender.AvatarURL,
		IsSystemMessage: false,
		Content:         msg.Content,
		CreatedAt:       msg.CreatedAt,
	}

	// Trigger Firebase Realtime Database push asynchronously
	firebasePath := fmt.Sprintf("chats/groups/%s/messages", gc.ID.String())
	utils.PushToFirebaseRealtime(firebasePath, response, s.cfg)

	return response, nil
}

func (s *conversationService) DeleteMessage(userID uuid.UUID, messageID uuid.UUID) error {
	// 1. Try to find message in Direct Messages
	msg, err := s.convRepo.GetMessageByID(messageID)
	if err == nil && msg != nil {
		if msg.SenderID != userID {
			return errors.New("unauthorized: you can only delete your own messages")
		}

		if err := s.convRepo.DeleteMessage(msg.ID); err != nil {
			return err
		}

		// Broadcast delete event to Firebase Realtime Database
		firebasePath := fmt.Sprintf("chats/direct/%s/messages", msg.ConversationID.String())
		utils.PushToFirebaseRealtime(firebasePath, map[string]interface{}{
			"event":      "message_deleted",
			"message_id": msg.ID.String(),
			"deleted_at": time.Now().UTC().Format(time.RFC3339),
		}, s.cfg)

		return nil
	}

	// 2. Try to find message in Group Messages
	gmsg, err := s.convRepo.GetGroupMessageByID(messageID)
	if err == nil && gmsg != nil {
		if gmsg.SenderID == nil || *gmsg.SenderID != userID {
			return errors.New("unauthorized: you can only delete your own messages")
		}

		if err := s.convRepo.DeleteGroupMessage(gmsg.ID); err != nil {
			return err
		}

		// Broadcast delete event to Firebase Realtime Database
		firebasePath := fmt.Sprintf("chats/groups/%s/messages", gmsg.GroupChatID.String())
		utils.PushToFirebaseRealtime(firebasePath, map[string]interface{}{
			"event":      "message_deleted",
			"message_id": gmsg.ID.String(),
			"deleted_at": time.Now().UTC().Format(time.RFC3339),
		}, s.cfg)

		return nil
	}

	return errors.New("message not found")
}
