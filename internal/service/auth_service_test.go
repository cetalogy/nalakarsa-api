package service

import (
	"errors"
	"testing"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"

	"github.com/google/uuid"
)

// MockUserRepository implements repository.UserRepository for testing
type mockUserRepository struct {
	users         map[string]*model.User
	refreshTokens map[string]*model.RefreshToken
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:         make(map[string]*model.User),
		refreshTokens: make(map[string]*model.RefreshToken),
	}
}

func (m *mockUserRepository) Create(user *model.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	user.ID = uuid.New()
	user.Profile.UserID = user.ID
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) UpdateProfile(profile *model.Profile) error {
	for _, u := range m.users {
		if u.ID == profile.UserID {
			u.Profile = *profile
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.Profile.AvatarURL = avatarURL
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) ListUsers(search, role string, page, limit int) ([]model.User, int64, error) {
	var result []model.User
	for _, u := range m.users {
		if role != "" && u.Role != role {
			continue
		}
		result = append(result, *u)
	}
	return result, int64(len(result)), nil
}

func (m *mockUserRepository) CreateRefreshToken(rt *model.RefreshToken) error {
	m.refreshTokens[rt.Token] = rt
	return nil
}

func (m *mockUserRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	rt, exists := m.refreshTokens[token]
	if !exists {
		return nil, nil
	}
	return rt, nil
}

func (m *mockUserRepository) DeleteRefreshToken(token string) error {
	delete(m.refreshTokens, token)
	return nil
}

func (m *mockUserRepository) DeleteRefreshTokensByUserID(userID uuid.UUID) error {
	for t, rt := range m.refreshTokens {
		if rt.UserID == userID {
			delete(m.refreshTokens, t)
		}
	}
	return nil
}

func (m *mockUserRepository) IncrementViewCount(userID uuid.UUID) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.Profile.ViewCount++
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) CountDiscussions(userID uuid.UUID) (int64, error) {
	return 0, nil
}

func TestRegisterAndLogin(t *testing.T) {
	mockRepo := newMockUserRepository()
	cfg := &config.Config{
		JWTSecret:            "test_secret",
		JWTRefreshSecret:     "test_refresh_secret",
		JWTAccessExpiration:  60,
		JWTRefreshExpiration: 3600,
	}

	authServ := NewAuthService(mockRepo, cfg)

	// Test case 1: Successful Registration
	regReq := dto.RegisterRequest{
		Email:          "test@nalakarsa.id",
		Password:       "password123",
		Role:           "profesional",
		NamaLengkap:    "Test User",
		Afiliasi:       "Nalakarsa Lab",
		Lokasi:         "Indonesia",
		BidangKeahlian: "Software",
	}

	regRes, err := authServ.Register(regReq)
	if err != nil {
		t.Fatalf("Expected registration to succeed, got error: %v", err)
	}

	if regRes.Email != regReq.Email {
		t.Errorf("Expected registered email %s, got %s", regReq.Email, regRes.Email)
	}

	// Test case 2: Duplicate Registration (Should fail)
	_, err = authServ.Register(regReq)
	if err == nil {
		t.Errorf("Expected duplicate registration to fail, but it succeeded")
	}

	// Test case 3: Successful Login
	loginReq := dto.LoginRequest{
		Email:    "test@nalakarsa.id",
		Password: "password123",
	}

	loginData, err := authServ.Login(loginReq)
	if err != nil {
		t.Fatalf("Expected login to succeed, got error: %v", err)
	}

	if loginData.AccessToken == "" {
		t.Errorf("Expected access token to be generated, got empty string")
	}

	if loginData.RefreshToken == "" {
		t.Errorf("Expected refresh token to be generated, got empty string")
	}

	if loginData.User.Email != loginReq.Email {
		t.Errorf("Expected logged-in user email %s, got %s", loginReq.Email, loginData.User.Email)
	}

	// Test case 4: Invalid Password Login (Should fail)
	invalidLoginReq := dto.LoginRequest{
		Email:    "test@nalakarsa.id",
		Password: "wrongpassword",
	}

	_, err = authServ.Login(invalidLoginReq)
	if err == nil {
		t.Errorf("Expected login with wrong password to fail, but it succeeded")
	}
}
