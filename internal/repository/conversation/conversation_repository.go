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
	ListMessages(conversationID uuid.UUID, limit int, cursor string) ([]model.Message, error)
	CountUnread(conversationID, userID uuid.UUID) (int64, error)
	CountTotalUnread(userID uuid.UUID) (int64, error)

	// Update conversation
	UpdateLastMessageAt(conversationID uuid.UUID) error
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
	err := r.db.Preload("Members.User.Profile").Where("id = ?", id).First(&conv).Error
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
	err := r.db.Preload("Members.User.Profile").
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
	err := r.db.Preload("Members.User.Profile").
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
	query := r.db.Preload("Sender.Profile").
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
