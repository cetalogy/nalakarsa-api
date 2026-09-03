package connectionservice

import (
	"errors"

	"nalakarsa/internal/common/constant"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	connectionrepository "nalakarsa/internal/repository/connection"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type ConnectionService interface {
	ListConnections(userID uuid.UUID, page, limit int) ([]dto.ConnectionResponse, int64, error)
	ListRequests(userID uuid.UUID, requestType string, page, limit int) ([]dto.ConnectionRequestResponse, int64, error)
	SendRequest(userID uuid.UUID, req dto.SendConnectionRequest) (uuid.UUID, error)
	AcceptRequest(userID uuid.UUID, requestID uuid.UUID) error
	RejectRequest(userID uuid.UUID, requestID uuid.UUID) error
	CancelRequest(userID uuid.UUID, requestID uuid.UUID) error
	RemoveConnection(userID uuid.UUID, targetUserID uuid.UUID) error
	GetSuggestions(userID uuid.UUID, limit int) ([]dto.UserSuggestionResponse, error)
}

type connectionService struct {
	connRepo connectionrepository.ConnectionRepository
	userRepo userrepository.UserRepository
}

func NewConnectionService(connRepo connectionrepository.ConnectionRepository, userRepo userrepository.UserRepository) ConnectionService {
	return &connectionService{connRepo: connRepo, userRepo: userRepo}
}

func (s *connectionService) ListConnections(userID uuid.UUID, page, limit int) ([]dto.ConnectionResponse, int64, error) {
	conns, total, err := s.connRepo.ListConnections(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.ConnectionResponse, len(conns))
	for i, c := range conns {
		other := c.Addressee
		if c.AddresseeID == userID {
			other = c.Requester
		}
		res[i] = dto.ConnectionResponse{
			ID:        c.ID,
			Status:    c.Status,
			CreatedAt: c.CreatedAt,
			User: dto.ConnectionUserResponse{
				ID:          other.ID,
				FullName:    other.FullName,
				Role:        other.Role,
				Affiliation: other.Affiliation,
				AvatarURL:   other.AvatarURL,
			},
		}
	}

	return res, total, nil
}

func (s *connectionService) ListRequests(userID uuid.UUID, requestType string, page, limit int) ([]dto.ConnectionRequestResponse, int64, error) {
	conns, total, err := s.connRepo.ListRequests(userID, requestType, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.ConnectionRequestResponse, len(conns))
	for i, c := range conns {
		other := c.Requester
		if requestType == "outgoing" {
			other = c.Addressee
		}
		res[i] = dto.ConnectionRequestResponse{
			ID:        c.ID,
			Status:    c.Status,
			CreatedAt: c.CreatedAt,
			User: dto.ConnectionUserResponse{
				ID:          other.ID,
				FullName:    other.FullName,
				Role:        other.Role,
				Affiliation: other.Affiliation,
				AvatarURL:   other.AvatarURL,
			},
		}
	}

	return res, total, nil
}

func (s *connectionService) SendRequest(userID uuid.UUID, req dto.SendConnectionRequest) (uuid.UUID, error) {
	if userID == req.TargetUserID {
		return uuid.Nil, errors.New("cannot send connection request to yourself")
	}
	target, err := s.userRepo.GetByID(req.TargetUserID)
	if err != nil {
		return uuid.Nil, err
	}
	if target == nil {
		return uuid.Nil, errors.New("target user not found")
	}
	existing, err := s.connRepo.GetByPair(userID, req.TargetUserID)
	if err != nil {
		return uuid.Nil, err
	}
	if existing != nil {
		if existing.Status == constant.ConnectionStatusAccepted {
			return uuid.Nil, errors.New("already connected with this user")
		}
		if existing.Status == constant.ConnectionStatusPending {
			return uuid.Nil, errors.New("connection request already exists")
		}
	}

	conn := &model.Connection{
		RequesterID: userID,
		AddresseeID: req.TargetUserID,
		Status:      constant.ConnectionStatusPending,
	}

	if err := s.connRepo.Create(conn); err != nil {
		return uuid.Nil, err
	}

	return conn.ID, nil
}

func (s *connectionService) AcceptRequest(userID uuid.UUID, requestID uuid.UUID) error {
	conn, err := s.connRepo.GetByID(requestID)
	if err != nil {
		return err
	}
	if conn == nil {
		return errors.New("connection request not found")
	}
	if conn.AddresseeID != userID {
		return errors.New("unauthorized to accept this request")
	}

	if conn.Status != constant.ConnectionStatusPending {
		return errors.New("connection request is no longer pending")
	}

	conn.Status = constant.ConnectionStatusAccepted
	return s.connRepo.Update(conn)
}

func (s *connectionService) RejectRequest(userID uuid.UUID, requestID uuid.UUID) error {
	conn, err := s.connRepo.GetByID(requestID)
	if err != nil {
		return err
	}
	if conn == nil {
		return errors.New("connection request not found")
	}
	if conn.AddresseeID != userID {
		return errors.New("unauthorized to reject this request")
	}

	if conn.Status != constant.ConnectionStatusPending {
		return errors.New("connection request is no longer pending")
	}

	conn.Status = constant.ConnectionStatusRejected
	return s.connRepo.Update(conn)
}

func (s *connectionService) CancelRequest(userID uuid.UUID, requestID uuid.UUID) error {
	conn, err := s.connRepo.GetByID(requestID)
	if err != nil {
		return err
	}
	if conn == nil {
		return errors.New("connection request not found")
	}
	if conn.RequesterID != userID {
		return errors.New("unauthorized to cancel this request")
	}

	if conn.Status != constant.ConnectionStatusPending {
		return errors.New("connection request is no longer pending")
	}

	return s.connRepo.Delete(requestID)
}

func (s *connectionService) RemoveConnection(userID uuid.UUID, targetUserID uuid.UUID) error {
	conn, err := s.connRepo.GetByPair(userID, targetUserID)
	if err != nil {
		return err
	}
	if conn == nil || conn.Status != constant.ConnectionStatusAccepted {
		return errors.New("connection not found")
	}

	return s.connRepo.Delete(conn.ID)
}

func (s *connectionService) GetSuggestions(userID uuid.UUID, limit int) ([]dto.UserSuggestionResponse, error) {
	users, _, err := s.userRepo.ListUsers("", "", 1, limit+10)
	if err != nil {
		return nil, err
	}

	var suggestions []dto.UserSuggestionResponse
	for _, u := range users {
		if u.ID == userID {
			continue
		}
		existing, _ := s.connRepo.GetByPair(userID, u.ID)
		if existing != nil {
			continue
		}

		mutualCount, _ := s.connRepo.CountMutual(userID, u.ID)

		suggestions = append(suggestions, dto.UserSuggestionResponse{
			ID:          u.ID,
			FullName:    u.FullName,
			Role:        u.Role,
			Affiliation: u.Affiliation,
			Expertise:   u.Expertise,
			AvatarURL:   u.AvatarURL,
			MutualCount: mutualCount,
		})

		if len(suggestions) >= limit {
			break
		}
	}

	return suggestions, nil
}
