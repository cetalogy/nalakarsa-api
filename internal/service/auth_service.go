package service

import (
	"errors"
	"strings"
	"time"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/model"
	"nalakarsa/internal/repository"
	"nalakarsa/internal/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginData, error)
	RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenData, error)
	Logout(token string) error
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
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

	// Create user & profile models
	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		SystemRole:   "user",
		Status:       "active",
		Profile: model.Profile{
			NamaLengkap:    req.NamaLengkap,
			GelarDepan:     req.GelarDepan,
			GelarBelakang:  req.GelarBelakang,
			Afiliasi:       req.Afiliasi,
			Lokasi:         req.Lokasi,
			BidangKeahlian: req.BidangKeahlian,
			Industry:       req.Industry,
			Bio:            req.Bio,
			Mission:        req.Mission,
			AvatarURL:      req.AvatarURL,
		},
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginData, error) {
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

	// Store Refresh Token in DB
	rtModel := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenPayload.Token,
		ExpiresAt: refreshTokenPayload.ExpiresAt,
	}
	if err := s.userRepo.CreateRefreshToken(rtModel); err != nil {
		return nil, err
	}

	return &dto.LoginData{
		AccessToken:  accessTokenPayload.Token,
		RefreshToken: refreshTokenPayload.Token,
		ExpiresIn:    s.cfg.JWTAccessExpiration,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *authService) RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenData, error) {
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

	// Delete old refresh token, save new one
	if err := s.userRepo.DeleteRefreshToken(req.RefreshToken); err != nil {
		return nil, err
	}

	newRt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenPayload.Token,
		ExpiresAt: refreshTokenPayload.ExpiresAt,
	}
	if err := s.userRepo.CreateRefreshToken(newRt); err != nil {
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
