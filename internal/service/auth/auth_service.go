package authservice

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	userrepository "nalakarsa/internal/repository/user"
	"nalakarsa/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(req dto.RegisterRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error)
	Login(req dto.LoginRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error)
	RefreshToken(req dto.RefreshTokenRequest, ctx *dto.AuthRequestContext) (*dto.RefreshTokenData, error)
	Logout(token string) error
	RequestPasswordReset(req dto.ForgotPasswordRequest) (string, error)
	ResetPassword(req dto.ResetPasswordRequest) error
}

type passwordResetClaims struct {
	Purpose string `json:"purpose"`
	UserID  string `json:"user_id"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo userrepository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo userrepository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(req dto.RegisterRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error) {
	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email is already registered")
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}
	securityAnswerHash, err := utils.HashPassword(normalizeSecurityAnswer(req.SecurityAnswer))
	if err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user model
	u := &model.User{
		Email:              req.Email,
		PasswordHash:       hashedPassword,
		SecurityQuestion:   req.SecurityQuestion,
		SecurityAnswerHash: securityAnswerHash,
		Role:               req.Role,
		SystemRole:         "user",
		Status:             "active",
		FirstName:          req.FirstName,
		MiddleName:         req.MiddleName,
		LastName:           req.LastName,
		FullName:           req.FullName,
		PrefixTitle:        req.PrefixTitle,
		SuffixTitle:        req.SuffixTitle,
		Affiliation:        req.Affiliation,
		Location:           req.Location,
		Expertise:          req.Expertise,
		Industry:           req.Industry,
		Bio:                req.Bio,
	}

	if err := s.userRepo.Create(u); err != nil {
		return nil, err
	}

	// Generate Tokens
	accessTokenPayload, err := utils.GenerateAccessToken(u.ID, u.Email, u.Role, s.cfg.JWTSecret, s.cfg.JWTAccessExpiration)
	if err != nil {
		return nil, err
	}
	refreshTokenPayload, err := utils.GenerateRefreshToken(u.ID, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshExpiration)
	if err != nil {
		return nil, err
	}

	if err := s.saveRefreshTokenWithSessionMeta(u.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, ctx); err != nil {
		return nil, err
	}

	return &dto.AuthData{
		Token:        accessTokenPayload.Token,
		AccessToken:  accessTokenPayload.Token,
		RefreshToken: refreshTokenPayload.Token,
		ExpiresIn:    s.cfg.JWTAccessExpiration,
		User: dto.UserResponse{
			ID:          u.ID,
			Email:       u.Email,
			Role:        u.Role,
			CreatedAt:   u.CreatedAt,
			FirstName:   u.FirstName,
			MiddleName:  u.MiddleName,
			LastName:    u.LastName,
			FullName:    u.FullName,
			PrefixTitle: u.PrefixTitle,
			SuffixTitle: u.SuffixTitle,
			Affiliation: u.Affiliation,
			Location:    u.Location,
			Expertise:   u.Expertise,
			Industry:    u.Industry,
			Bio:         u.Bio,
			AvatarURL:   u.AvatarURL,
		},
	}, nil
}

func validatePassword(password string) error {
	if len([]rune(password)) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len([]rune(password)) > 72 {
		return errors.New("password must not exceed 72 characters")
	}
	var upper, lower, number, special bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsLower(char):
			lower = true
		case unicode.IsNumber(char):
			number = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		}
	}
	if !upper || !lower || !number || !special {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}
	return nil
}

func (s *authService) Login(req dto.LoginRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check account status
	if u.Status != "active" {
		return nil, errors.New("account is suspended")
	}

	// Verify password
	if !utils.ComparePassword(u.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate Access Token
	accessTokenPayload, err := utils.GenerateAccessToken(
		u.ID,
		u.Email,
		u.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTAccessExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshTokenPayload, err := utils.GenerateRefreshToken(
		u.ID,
		s.cfg.JWTRefreshSecret,
		s.cfg.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	if err := s.saveRefreshTokenWithSessionMeta(u.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, ctx); err != nil {
		return nil, err
	}

	return &dto.AuthData{
		Token:        accessTokenPayload.Token,
		AccessToken:  accessTokenPayload.Token,
		RefreshToken: refreshTokenPayload.Token,
		ExpiresIn:    s.cfg.JWTAccessExpiration,
		User: dto.UserResponse{
			ID:          u.ID,
			Email:       u.Email,
			Role:        u.Role,
			CreatedAt:   u.CreatedAt,
			FirstName:   u.FirstName,
			MiddleName:  u.MiddleName,
			LastName:    u.LastName,
			FullName:    u.FullName,
			PrefixTitle: u.PrefixTitle,
			SuffixTitle: u.SuffixTitle,
			Affiliation: u.Affiliation,
			Location:    u.Location,
			Expertise:   u.Expertise,
			Industry:    u.Industry,
			Bio:         u.Bio,
			AvatarURL:   u.AvatarURL,
		},
	}, nil
}

func (s *authService) RefreshToken(req dto.RefreshTokenRequest, ctx *dto.AuthRequestContext) (*dto.RefreshTokenData, error) {
	// Find token in database
	rt, err := s.userRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check expiration
	if time.Now().After(rt.ExpiresAt) {
		_ = s.userRepo.DeleteRefreshToken(req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	// Get User details
	u, err := s.userRepo.GetByID(rt.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Generate new Access Token
	accessTokenPayload, err := utils.GenerateAccessToken(
		u.ID,
		u.Email,
		u.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTAccessExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Generate new Refresh Token (Refresh Token Rotation)
	refreshTokenPayload, err := utils.GenerateRefreshToken(
		u.ID,
		s.cfg.JWTRefreshSecret,
		s.cfg.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Save new token first so we don't silently drop sessions on storage failure.
	if err := s.saveRefreshTokenWithSessionMeta(u.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, s.withFallbackSessionContext(ctx, rt)); err != nil {
		return nil, err
	}

	// Revoke old refresh token after the new token is safely persisted.
	if err := s.userRepo.DeleteRefreshToken(req.RefreshToken); err != nil {
		// Best-effort rollback: remove the newly issued token so we don't keep extra active tokens.
		if rollbackErr := s.userRepo.DeleteRefreshToken(refreshTokenPayload.Token); rollbackErr != nil {
			return nil, fmt.Errorf("failed to revoke old refresh token and rollback new token: %w", rollbackErr)
		}
		return nil, err
	}

	return &dto.RefreshTokenData{
		AccessToken:  accessTokenPayload.Token,
		RefreshToken: refreshTokenPayload.Token,
		ExpiresIn:    s.cfg.JWTAccessExpiration,
	}, nil
}

func (s *authService) Logout(token string) error {
	return s.userRepo.DeleteRefreshToken(token)
}

func (s *authService) RequestPasswordReset(req dto.ForgotPasswordRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return "", err
	}
	if user == nil || user.SecurityQuestion != req.SecurityQuestion || !utils.ComparePassword(user.SecurityAnswerHash, normalizeSecurityAnswer(req.SecurityAnswer)) {
		return "", errors.New("security question answer is incorrect")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, passwordResetClaims{
		Purpose: "password_reset",
		UserID:  user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) ResetPassword(req dto.ResetPasswordRequest) error {
	claims := &passwordResetClaims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid reset token signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid || claims.Purpose != "password_reset" {
		return errors.New("invalid or expired password reset token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return errors.New("invalid password reset user")
	}
	if err := validatePassword(req.NewPassword); err != nil {
		return err
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash
	return s.userRepo.UpdateProfile(user)
}

func normalizeSecurityAnswer(answer string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(answer)), " "))
}

func (s *authService) saveRefreshTokenWithSessionMeta(
	userID uuid.UUID,
	token string,
	expiresAt time.Time,
	ctx *dto.AuthRequestContext,
) error {
	if s.cfg.MaxActiveRefreshTokens > 0 {
		if err := s.limitActiveRefreshTokens(userID, s.cfg.MaxActiveRefreshTokens-1); err != nil {
			return err
		}
	}

	deviceInfo := "unknown-device"
	userAgent := ""
	ipAddress := ""
	if ctx != nil {
		if strings.TrimSpace(ctx.DeviceInfo) != "" {
			deviceInfo = strings.TrimSpace(ctx.DeviceInfo)
		}
		userAgent = strings.TrimSpace(ctx.UserAgent)
		ipAddress = strings.TrimSpace(ctx.IPAddress)
	}

	refreshToken := &model.RefreshToken{
		UserID:     userID,
		Token:      token,
		ExpiresAt:  expiresAt,
		DeviceInfo: deviceInfo,
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
	}

	return s.userRepo.CreateRefreshToken(refreshToken)
}

func (s *authService) withFallbackSessionContext(
	ctx *dto.AuthRequestContext,
	legacy *model.RefreshToken,
) *dto.AuthRequestContext {
	fallback := &dto.AuthRequestContext{}
	if ctx != nil {
		*fallback = *ctx
	}

	if legacy == nil {
		return fallback
	}

	if strings.TrimSpace(fallback.DeviceInfo) == "" {
		fallback.DeviceInfo = legacy.DeviceInfo
	}
	if strings.TrimSpace(fallback.UserAgent) == "" {
		fallback.UserAgent = legacy.UserAgent
	}
	if strings.TrimSpace(fallback.IPAddress) == "" {
		fallback.IPAddress = legacy.IPAddress
	}

	return fallback
}

func (s *authService) limitActiveRefreshTokens(userID uuid.UUID, keepLatest int) error {
	if keepLatest < 0 {
		return nil
	}

	count, err := s.userRepo.CountActiveRefreshTokens(userID)
	if err != nil {
		return err
	}

	if int(count) <= keepLatest {
		return nil
	}

	return s.userRepo.DeleteOldestRefreshTokensByUser(userID, keepLatest)
}
