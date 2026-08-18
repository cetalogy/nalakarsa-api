package conversationrepository

import (
	"errors"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(conv *model.Conversation) error
	GetByID(id uuid.UUID) (*model.Conversation, error)
	GetDirectByPair(userA, userB uuid.UUID) (*model.Conversation, error)
	ListByUser(userID uuid.UUID, page, limit int) ([]model.Conversation, int64, error)

	// Members
	AddMember(member *model.ConversationMember) error
	GetMember(conversationID, userID uuid.UUID) (*model.ConversationMember, error)
	UpdateLastRead(conversationID, userID, messageID uuid.UUID) error

	// Messages
	CreateMessage(msg *model.Message) error
	GetMessageByID(id uuid.UUID) (*model.Message, error)
	DeleteMessage(id uuid.UUID) error
	ListMessages(conversationID uuid.UUID, limit int, cursor string) ([]model.Message, error)
	CountUnread(conversationID, userID uuid.UUID) (int64, error)
	CountTotalUnread(userID uuid.UUID) (int64, error)

	// Update conversation
	UpdateLastMessageAt(conversationID uuid.UUID) error

	// Group Chats (FE Contract Specification)
	ListGroupChatsByUser(userID uuid.UUID) ([]model.GroupChat, error)
	GetGroupChatByID(id uuid.UUID) (*model.GroupChat, error)
	IsGroupChatMember(groupChatID, userID uuid.UUID) (bool, error)
	ListGroupMessages(groupChatID uuid.UUID, page, limit int) ([]model.GroupMessage, int64, error)
	CreateGroupMessage(msg *model.GroupMessage) error
	GetGroupMessageByID(id uuid.UUID) (*model.GroupMessage, error)
	DeleteGroupMessage(id uuid.UUID) error
	UpdateGroupChatLastMessage(groupChatID uuid.UUID, lastMessage string) error
}

type pgConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &pgConversationRepository{db: db}
}

func (r *pgConversationRepository) Create(conv *model.Conversation) error {
	return r.db.Create(conv).Error
}

func (r *pgConversationRepository) GetByID(id uuid.UUID) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.Preload("Members.User").Where("id = ?", id).First(&conv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conv, nil
}

func (r *pgConversationRepository) GetDirectByPair(userA, userB uuid.UUID) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.Preload("Members.User").
		Where("type = 'direct' AND id IN ("+
			"SELECT cm1.conversation_id FROM conversation_members cm1 "+
			"INNER JOIN conversation_members cm2 ON cm1.conversation_id = cm2.conversation_id "+
			"WHERE cm1.user_id = ? AND cm2.user_id = ?"+
			")", userA, userB).
		First(&conv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conv, nil
}

func (r *pgConversationRepository) ListByUser(userID uuid.UUID, page, limit int) ([]model.Conversation, int64, error) {
	var convs []model.Conversation
	var total int64

	query := r.db.Model(&model.Conversation{}).
		Where("id IN (SELECT conversation_id FROM conversation_members WHERE user_id = ?)", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("Members.User").
		Where("id IN (SELECT conversation_id FROM conversation_members WHERE user_id = ?)", userID).
		Order("COALESCE(last_message_at, created_at) DESC").
		Limit(limit).Offset(offset).Find(&convs).Error
	return convs, total, err
}

// --- Members ---

func (r *pgConversationRepository) AddMember(member *model.ConversationMember) error {
	return r.db.Create(member).Error
}

func (r *pgConversationRepository) GetMember(conversationID, userID uuid.UUID) (*model.ConversationMember, error) {
	var member model.ConversationMember
	err := r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *pgConversationRepository) UpdateLastRead(conversationID, userID, messageID uuid.UUID) error {
	return r.db.Model(&model.ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_message_id", messageID).Error
}

// --- Messages ---

func (r *pgConversationRepository) CreateMessage(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *pgConversationRepository) ListMessages(conversationID uuid.UUID, limit int, cursor string) ([]model.Message, error) {
	var messages []model.Message
	query := r.db.Preload("Sender").
		Where("conversation_id = ?", conversationID)

	if cursor != "" {
		cursorID, err := uuid.Parse(cursor)
		if err == nil {
			// Get the created_at of the cursor message
			var cursorMsg model.Message
			if err := r.db.Select("created_at").Where("id = ?", cursorID).First(&cursorMsg).Error; err == nil {
				query = query.Where("messages.created_at < ?", cursorMsg.CreatedAt)
			}
		}
	}

	err := query.Order("messages.created_at desc").Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *pgConversationRepository) CountUnread(conversationID, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM messages m
		WHERE m.conversation_id = ?
		AND m.sender_id != ?
		AND m.id > COALESCE(
			(SELECT last_read_message_id FROM conversation_members WHERE conversation_id = ? AND user_id = ?),
			'00000000-0000-0000-0000-000000000000'::uuid
		)
	`, conversationID, userID, conversationID, userID).Scan(&count).Error
	return count, err
}

func (r *pgConversationRepository) CountTotalUnread(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM messages m
		INNER JOIN conversation_members cm ON cm.conversation_id = m.conversation_id AND cm.user_id = ?
		WHERE m.sender_id != ?
		AND m.id > COALESCE(cm.last_read_message_id, '00000000-0000-0000-0000-000000000000'::uuid)
	`, userID, userID).Scan(&count).Error
	return count, err
}

func (r *pgConversationRepository) UpdateLastMessageAt(conversationID uuid.UUID) error {
	return r.db.Exec(
		"UPDATE conversations SET last_message_at = NOW(), updated_at = NOW() WHERE id = ?",
		conversationID,
	).Error
}

func (r *pgConversationRepository) ListGroupChatsByUser(userID uuid.UUID) ([]model.GroupChat, error) {
	var groupChats []model.GroupChat
	subQuery := r.db.Model(&model.GroupChatMember{}).Select("group_chat_id").Where("user_id = ?", userID)

	err := r.db.Preload("Members.User").
		Where("id IN (?)", subQuery).
		Order("created_at desc").
		Find(&groupChats).Error

	return groupChats, err
}

func (r *pgConversationRepository) GetGroupChatByID(id uuid.UUID) (*model.GroupChat, error) {
	var groupChat model.GroupChat
	err := r.db.Preload("Members.User").
		Where("id = ? OR topic_id = ? OR project_id = ?", id, id, id).
		First(&groupChat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &groupChat, nil
}

func (r *pgConversationRepository) IsGroupChatMember(groupChatID, userID uuid.UUID) (bool, error) {
	// First resolve the real groupChat.ID in case groupChatID was a topic_id or project_id
	gc, err := r.GetGroupChatByID(groupChatID)
	if err != nil || gc == nil {
		return false, err
	}

	var count int64
	err = r.db.Model(&model.GroupChatMember{}).
		Where("group_chat_id = ? AND user_id = ?", gc.ID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *pgConversationRepository) ListGroupMessages(groupChatID uuid.UUID, page, limit int) ([]model.GroupMessage, int64, error) {
	var messages []model.GroupMessage
	var total int64

	query := r.db.Model(&model.GroupMessage{}).Where("group_chat_id = ?", groupChatID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Sender").
		Order("created_at asc").
		Limit(limit).Offset(offset).
		Find(&messages).Error

	return messages, total, err
}

func (r *pgConversationRepository) CreateGroupMessage(msg *model.GroupMessage) error {
	return r.db.Create(msg).Error
}

func (r *pgConversationRepository) UpdateGroupChatLastMessage(groupChatID uuid.UUID, lastMessage string) error {
	return r.db.Model(&model.GroupChat{}).
		Where("id = ?", groupChatID).
		Updates(map[string]interface{}{
			"last_message":      lastMessage,
			"last_message_time": gorm.Expr("NOW()"),
		}).Error
}

func (r *pgConversationRepository) GetMessageByID(id uuid.UUID) (*model.Message, error) {
	var msg model.Message
	err := r.db.Where("id = ?", id).First(&msg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

func (r *pgConversationRepository) DeleteMessage(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Message{}).Error
}

func (r *pgConversationRepository) GetGroupMessageByID(id uuid.UUID) (*model.GroupMessage, error) {
	var msg model.GroupMessage
	err := r.db.Where("id = ?", id).First(&msg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

func (r *pgConversationRepository) DeleteGroupMessage(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.GroupMessage{}).Error
}
