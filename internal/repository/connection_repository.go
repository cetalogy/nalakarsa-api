package repository

import (
	"errors"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConnectionRepository interface {
	Create(conn *model.Connection) error
	GetByID(id uuid.UUID) (*model.Connection, error)
	GetByPair(userA, userB uuid.UUID) (*model.Connection, error)
	ListConnections(userID uuid.UUID, page, limit int) ([]model.Connection, int64, error)
	ListRequests(userID uuid.UUID, requestType string, page, limit int) ([]model.Connection, int64, error)
	Update(conn *model.Connection) error
	Delete(id uuid.UUID) error
	CountAccepted(userID uuid.UUID) (int64, error)
	CountMutual(userA, userB uuid.UUID) (int64, error)
}

type pgConnectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) ConnectionRepository {
	return &pgConnectionRepository{db: db}
}

func (r *pgConnectionRepository) Create(conn *model.Connection) error {
	return r.db.Create(conn).Error
}

func (r *pgConnectionRepository) GetByID(id uuid.UUID) (*model.Connection, error) {
	var conn model.Connection
	err := r.db.Preload("Requester.Profile").Preload("Addressee.Profile").
		Where("id = ?", id).First(&conn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conn, nil
}

func (r *pgConnectionRepository) GetByPair(userA, userB uuid.UUID) (*model.Connection, error) {
	var conn model.Connection
	err := r.db.Where(
		"(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)",
		userA, userB, userB, userA,
	).First(&conn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conn, nil
}

func (r *pgConnectionRepository) ListConnections(userID uuid.UUID, page, limit int) ([]model.Connection, int64, error) {
	var conns []model.Connection
	var total int64

	query := r.db.Model(&model.Connection{}).
		Preload("Requester.Profile").Preload("Addressee.Profile").
		Where("status = 'accepted' AND (requester_id = ? OR addressee_id = ?)", userID, userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("updated_at desc").Find(&conns).Error
	return conns, total, err
}

func (r *pgConnectionRepository) ListRequests(userID uuid.UUID, requestType string, page, limit int) ([]model.Connection, int64, error) {
	var conns []model.Connection
	var total int64

	query := r.db.Model(&model.Connection{}).
		Preload("Requester.Profile").Preload("Addressee.Profile").
		Where("status = 'pending'")

	if requestType == "incoming" {
		query = query.Where("addressee_id = ?", userID)
	} else {
		query = query.Where("requester_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&conns).Error
	return conns, total, err
}

func (r *pgConnectionRepository) Update(conn *model.Connection) error {
	return r.db.Save(conn).Error
}

func (r *pgConnectionRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Connection{}).Error
}

func (r *pgConnectionRepository) CountAccepted(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Connection{}).
		Where("status = 'accepted' AND (requester_id = ? OR addressee_id = ?)", userID, userID).
		Count(&count).Error
	return count, err
}

func (r *pgConnectionRepository) CountMutual(userA, userB uuid.UUID) (int64, error) {
	var count int64
	// Mutual connections: users connected to both userA and userB
	err := r.db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT CASE WHEN requester_id = ? THEN addressee_id ELSE requester_id END as uid
			FROM connections WHERE status = 'accepted' AND (requester_id = ? OR addressee_id = ?)
			INTERSECT
			SELECT CASE WHEN requester_id = ? THEN addressee_id ELSE requester_id END as uid
			FROM connections WHERE status = 'accepted' AND (requester_id = ? OR addressee_id = ?)
		) AS mutual
	`, userA, userA, userA, userB, userB, userB).Scan(&count).Error
	return count, err
}
