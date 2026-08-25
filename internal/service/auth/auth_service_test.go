package authservice

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

func (m *mockUserRepository) Create(u *model.User) error {
	if _, exists := m.users[u.Email]; exists {
		return errors.New("user already exists")
	}
	u.ID = uuid.New()
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	u, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) GetByIDOrIdentifier(identifier string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == identifier || u.FullName == identifier {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) UpdateProfile(u *model.User) error {
	for _, existing := range m.users {
		if existing.ID == u.ID {
			existing.FirstName = u.FirstName
			existing.MiddleName = u.MiddleName
			existing.LastName = u.LastName
			existing.FullName = u.FullName
			existing.PrefixTitle = u.PrefixTitle
			existing.SuffixTitle = u.SuffixTitle
			existing.Affiliation = u.Affiliation
			existing.Location = u.Location
			existing.Expertise = u.Expertise
			existing.Industry = u.Industry
			existing.Bio = u.Bio
			existing.Mission = u.Mission
			existing.AvatarURL = u.AvatarURL
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.AvatarURL = avatarURL
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

func (m *mockUserRepository) CountActiveRefreshTokens(userID uuid.UUID) (int64, error) {
	var count int64
	for _, rt := range m.refreshTokens {
		if rt.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockUserRepository) DeleteOldestRefreshTokensByUser(userID uuid.UUID, keepLatest int) error {
	if keepLatest < 0 {
		keepLatest = 0
	}

	// simple strategy: remove all first, then keep the newest by insertion not guaranteed in mock
	var tokens []*model.RefreshToken
	for t, rt := range m.refreshTokens {
		if rt.UserID == userID {
			_ = t
			tokens = append(tokens, rt)
		}
	}

	if len(tokens) <= keepLatest {
		return nil
	}

	removeCount := len(tokens) - keepLatest
	removed := 0
	for token, rt := range m.refreshTokens {
		if rt.UserID == userID && removed < removeCount {
			delete(m.refreshTokens, token)
			removed++
		}
	}

	return nil
}

func (m *mockUserRepository) IncrementViewCount(userID uuid.UUID) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.ViewCount++
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) CountDiscussions(userID uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockUserRepository) ToggleFollow(followerID, followingID uuid.UUID) (bool, error) {
	return true, nil
}

func (m *mockUserRepository) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *mockUserRepository) GetFollowers(targetUserID uuid.UUID, page, limit int) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}

func (m *mockUserRepository) GetFollowing(targetUserID uuid.UUID, page, limit int) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}

func (m *mockUserRepository) CountFollowers(targetUserID uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockUserRepository) CountFollowing(targetUserID uuid.UUID) (int64, error) {
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
		Email:       "test@nalakarsa.id",
		Password:    "Password123!",
		Role:        "profesional",
		FirstName:   "Test",
		LastName:    "User",
		FullName:    "Test User",
		Affiliation: "Nalakarsa Lab",
		Location:    "Indonesia",
		Expertise:   "Software",
	}

	regRes, err := authServ.Register(regReq, nil)
	if err != nil {
		t.Fatalf("Expected registration to succeed, got error: %v", err)
	}

	if regRes.User.Email != regReq.Email {
		t.Errorf("Expected registered email %s, got %s", regReq.Email, regRes.User.Email)
	}

	// Test case 2: Duplicate Registration (Should fail)
	_, err = authServ.Register(regReq, nil)
	if err == nil {
		t.Errorf("Expected duplicate registration to fail, but it succeeded")
	}

	// Test case 3: Successful Login
	loginReq := dto.LoginRequest{
		Email:    "test@nalakarsa.id",
		Password: "Password123!",
	}

	loginData, err := authServ.Login(loginReq, nil)
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

	_, err = authServ.Login(invalidLoginReq, nil)
	if err == nil {
		t.Errorf("Expected login with wrong password to fail, but it succeeded")
	}
}
