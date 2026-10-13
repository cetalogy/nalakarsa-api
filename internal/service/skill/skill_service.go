package skill

import (
	"errors"
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/reference"
	skillrepo "nalakarsa/internal/repository/skill"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Search(search string, limit int) ([]dto.SkillResponse, error)
	Create(name string) (dto.SkillResponse, error)
	ListByUser(userID uuid.UUID) ([]dto.SkillResponse, error)
	AddToUser(userID uuid.UUID, skillID *uuid.UUID, name *string) (dto.SkillResponse, error)
	RemoveFromUser(userID, skillID uuid.UUID) error
}

type service struct{ repo skillrepo.Repository }

func NewSkillService(repo skillrepo.Repository) Service { return &service{repo: repo} }

func (s *service) Search(search string, limit int) ([]dto.SkillResponse, error) {
	items, err := s.repo.Search(search, limit)
	return mapSkills(items), err
}

func (s *service) Create(name string) (dto.SkillResponse, error) {
	normalizedName, err := reference.NormalizeName(name)
	if err != nil {
		return dto.SkillResponse{}, err
	}
	item, err := s.repo.Create(&model.Skill{Name: normalizedName, IsActive: true})
	if err != nil {
		return dto.SkillResponse{}, err
	}
	return toResponse(*item), nil
}

func (s *service) ListByUser(userID uuid.UUID) ([]dto.SkillResponse, error) {
	items, err := s.repo.ListByUser(userID)
	return mapSkills(items), err
}

func (s *service) AddToUser(userID uuid.UUID, skillID *uuid.UUID, name *string) (dto.SkillResponse, error) {
	if skillID == nil && (name == nil || strings.TrimSpace(*name) == "") {
		return dto.SkillResponse{}, errors.New("skill_id or name is required")
	}
	var item *model.Skill
	var err error
	if skillID != nil {
		item, err = s.repo.GetByID(*skillID)
	} else {
		var result dto.SkillResponse
		result, err = s.Create(*name)
		if err != nil {
			return dto.SkillResponse{}, err
		}
		item = &model.Skill{ID: result.ID, Name: result.Name}
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SkillResponse{}, errors.New("skill not found")
		}
		return dto.SkillResponse{}, err
	}
	if err := s.repo.AddToUser(userID, item.ID); err != nil {
		return dto.SkillResponse{}, err
	}
	return toResponse(*item), nil
}

func (s *service) RemoveFromUser(userID, skillID uuid.UUID) error {
	if err := s.repo.RemoveFromUser(userID, skillID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("skill is not assigned to this user")
		}
		return err
	}
	return nil
}

func mapSkills(items []model.Skill) []dto.SkillResponse {
	result := make([]dto.SkillResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toResponse(item))
	}
	return result
}

func toResponse(item model.Skill) dto.SkillResponse {
	return dto.SkillResponse{ID: item.ID, Name: item.Name}
}
