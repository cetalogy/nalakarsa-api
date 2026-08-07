package authservice

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/utils"
	userrepository "nalakarsa/internal/repository/user"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req dto.RegisterRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error)
	Login(req dto.LoginRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error)
	RefreshToken(req dto.RefreshTokenRequest, ctx *dto.AuthRequestContext) (*dto.RefreshTokenData, error)
	Logout(token string) error
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

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user model
	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		SystemRole:   "user",
		Status:       "active",
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		LastName:     req.LastName,
		FullName:     req.FullName,
		PrefixTitle:  req.PrefixTitle,
		SuffixTitle:  req.SuffixTitle,
		Affiliation:  req.Affiliation,
		Location:     req.Location,
		Expertise:    req.Expertise,
		Industry:     req.Industry,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate Tokens
	accessTokenPayload, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.cfg.JWTSecret, s.cfg.JWTAccessExpiration)
	if err != nil {
		return nil, err
	}
	refreshTokenPayload, err := utils.GenerateRefreshToken(user.ID, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshExpiration)
	if err != nil {
		return nil, err
	}

	if err := s.saveRefreshTokenWithSessionMeta(user.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, ctx); err != nil {
		return nil, err
	}

	return &dto.AuthData{
		Token:        accessTokenPayload.Token,
		AccessToken:  accessTokenPayload.Token,
		RefreshToken: refreshTokenPayload.Token,
		ExpiresIn:    s.cfg.JWTAccessExpiration,
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			FirstName:   user.FirstName,
			MiddleName:  user.MiddleName,
			LastName:    user.LastName,
			FullName:    user.FullName,
			PrefixTitle: user.PrefixTitle,
			SuffixTitle: user.SuffixTitle,
			Affiliation: user.Affiliation,
			Location:    user.Location,
			Expertise:   user.Expertise,
			Industry:    user.Industry,
			Mission:     user.Mission,
			AvatarURL:   user.AvatarURL,
		},
	}, nil
}

func (s *authService) Login(req dto.LoginRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check account status
	if user.Status != "active" {
		return nil, errors.New("account is suspended")
	}

	// Verify password
	if !utils.ComparePassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate Access Token
	accessTokenPayload, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTAccessExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshTokenPayload, err := utils.GenerateRefreshToken(
		user.ID,
		s.cfg.JWTRefreshSecret,
		s.cfg.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	if err := s.saveRefreshTokenWithSessionMeta(user.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, ctx); err != nil {
		return nil, err
	}

	return &dto.AuthData{
		Token:        accessTokenPayload.Token,
		AccessToken:  accessTokenPayload.Token,
		RefreshToken: refreshTokenPayload.Token,
		ExpiresIn:    s.cfg.JWTAccessExpiration,
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			FirstName:   user.FirstName,
			MiddleName:  user.MiddleName,
			LastName:    user.LastName,
			FullName:    user.FullName,
			PrefixTitle: user.PrefixTitle,
			SuffixTitle: user.SuffixTitle,
			Affiliation: user.Affiliation,
			Location:    user.Location,
			Expertise:   user.Expertise,
			Industry:    user.Industry,
			Mission:     user.Mission,
			AvatarURL:   user.AvatarURL,
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
	user, err := s.userRepo.GetByID(rt.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new Access Token
	accessTokenPayload, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTAccessExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Generate new Refresh Token (Refresh Token Rotation)
	refreshTokenPayload, err := utils.GenerateRefreshToken(
		user.ID,
		s.cfg.JWTRefreshSecret,
		s.cfg.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Save new token first so we don't silently drop sessions on storage failure.
	if err := s.saveRefreshTokenWithSessionMeta(user.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, s.withFallbackSessionContext(ctx, rt)); err != nil {
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
