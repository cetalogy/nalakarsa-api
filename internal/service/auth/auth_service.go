package authservice

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"strings"
	"time"
	"unicode"

	"nalakarsa/internal/config"
	"nalakarsa/internal/contentfilter"
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
	VerifyEmail(req dto.VerifyEmailRequest) error
	ResendVerificationEmail(req dto.ResendVerificationEmailRequest) error
}

type passwordResetClaims struct {
	Purpose string `json:"purpose"`
	UserID  string `json:"user_id"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo userrepository.UserRepository
	cfg      *config.Config
	email    *emailSender
}

func NewAuthService(userRepo userrepository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg, email: newEmailSender(cfg)}
}

func (s *authService) Register(req dto.RegisterRequest, ctx *dto.AuthRequestContext) (*dto.AuthData, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
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
	values := []string{req.FirstName, req.LastName, req.FullName, req.Affiliation, req.Location, req.Expertise, req.Industry, req.Bio}
	if req.MiddleName != nil {
		values = append(values, *req.MiddleName)
	}
	for _, value := range values {
		if err := contentfilter.Validate(value); err != nil {
			return nil, err
		}
	}
	securityAnswerHash, err := utils.HashPassword(normalizeSecurityAnswer(req.SecurityAnswer))
	if err != nil {
		return nil, err
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
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
	if s.email != nil {
		if err := s.sendVerificationEmail(u); err != nil {
			logEmailStatus("VERIFICATION", u.Email, err)
		} else {
			logEmailStatus("VERIFICATION", u.Email, nil)
		}
	}
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
	if u.Status != "active" {
		return nil, errors.New("account is suspended")
	}
	if !utils.ComparePassword(u.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}
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
	rt, err := s.userRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, errors.New("invalid refresh token")
	}
	if time.Now().After(rt.ExpiresAt) {
		_ = s.userRepo.DeleteRefreshToken(req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}
	u, err := s.userRepo.GetByID(rt.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}
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
	refreshTokenPayload, err := utils.GenerateRefreshToken(
		u.ID,
		s.cfg.JWTRefreshSecret,
		s.cfg.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, err
	}
	if err := s.saveRefreshTokenWithSessionMeta(u.ID, refreshTokenPayload.Token, refreshTokenPayload.ExpiresAt, s.withFallbackSessionContext(ctx, rt)); err != nil {
		return nil, err
	}
	if err := s.userRepo.DeleteRefreshToken(req.RefreshToken); err != nil {
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
	if s.email != nil {
		resetLink := strings.TrimRight(s.cfg.FrontendURL, "/") + "/reset-password?token=" + signedToken
		if err := s.email.sendPasswordReset(user.Email, user.FullName, resetLink); err != nil {
			logEmailStatus("PASSWORD RESET", user.Email, err)
			return "", fmt.Errorf("failed to send password reset email: %w", err)
		}
		logEmailStatus("PASSWORD RESET", user.Email, nil)
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
	if err := s.userRepo.UpdateProfile(user); err != nil {
		return err
	}
	if err := s.userRepo.DeleteRefreshTokensByUserID(user.ID); err != nil {
		return err
	}
	if s.email != nil {
		if err := s.email.sendPasswordChanged(user.Email, user.FullName); err != nil {
			logEmailStatus("PASSWORD CHANGED", user.Email, err)
		} else {
			logEmailStatus("PASSWORD CHANGED", user.Email, nil)
		}
	}
	return nil
}

func (s *authService) VerifyEmail(req dto.VerifyEmailRequest) error {
	user, err := s.userRepo.GetByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil || user == nil {
		return errors.New("email atau kode verifikasi tidak valid")
	}
	if user.EmailVerifiedAt != nil || user.EmailVerificationExpiresAt == nil || time.Now().After(*user.EmailVerificationExpiresAt) {
		return errors.New("kode verifikasi tidak valid atau sudah kedaluwarsa")
	}
	if hashVerificationCode(req.Code) != user.EmailVerificationCodeHash {
		return errors.New("email atau kode verifikasi tidak valid")
	}
	if err := s.userRepo.MarkEmailVerified(user.ID); err != nil {
		return err
	}
	if s.email != nil {
		loginLink := strings.TrimRight(s.cfg.FrontendURL, "/") + "/login"
		if err := s.email.sendWelcome(user.Email, user.FullName, loginLink); err != nil {
			logEmailStatus("WELCOME", user.Email, err)
		} else {
			logEmailStatus("WELCOME", user.Email, nil)
		}
	}
	return nil
}

func (s *authService) ResendVerificationEmail(req dto.ResendVerificationEmailRequest) error {
	user, err := s.userRepo.GetByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil || user == nil || user.EmailVerifiedAt != nil {
		return nil
	}
	err = s.sendVerificationEmail(user)
	logEmailStatus("VERIFICATION RESEND", user.Email, err)
	return err
}

func logEmailStatus(emailType, recipient string, err error) {
	if err != nil {
		log.Printf("[EMAIL][%s][FAILED] recipient=%s error=%v", emailType, recipient, err)
		return
	}
	log.Printf("[EMAIL][%s][SENT] recipient=%s", emailType, recipient)
}

func (s *authService) sendVerificationEmail(user *model.User) error {
	code, err := generateVerificationCode()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.userRepo.SetEmailVerification(user.ID, hashVerificationCode(code), expiresAt); err != nil {
		return err
	}
	if s.email == nil {
		return nil
	}
	return s.email.sendVerification(user.Email, user.FullName, code)
}

func generateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func hashVerificationCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
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

type emailSender struct{ cfg *config.Config }

func newEmailSender(cfg *config.Config) *emailSender { return &emailSender{cfg: cfg} }

func (s *emailSender) sendWelcome(to, name, loginLink string) error {
	message, err := buildWelcomeEmail(WelcomeData{RecipientName: name, LoginLink: loginLink})
	if err != nil {
		return err
	}
	return s.send(to, message)
}

func (s *emailSender) sendVerification(to, name, code string) error {
	message, err := buildVerificationEmail(VerificationData{RecipientName: name, Code: code})
	if err != nil {
		return err
	}
	return s.send(to, message)
}

func (s *emailSender) sendPasswordReset(to, name, resetLink string) error {
	message, err := buildPasswordResetEmail(PasswordResetData{RecipientName: name, ResetLink: resetLink})
	if err != nil {
		return err
	}
	return s.send(to, message)
}

func (s *emailSender) sendPasswordChanged(to, name string) error {
	message, err := buildPasswordChangedEmail(PasswordChangedData{RecipientName: name})
	if err != nil {
		return err
	}
	return s.send(to, message)
}

func (s *emailSender) send(to string, message emailMessage) error {
	if strings.TrimSpace(s.cfg.SMTPHost) == "" || strings.TrimSpace(s.cfg.SMTPUsername) == "" || strings.TrimSpace(s.cfg.SMTPPassword) == "" || strings.TrimSpace(s.cfg.SMTPFrom) == "" {
		return fmt.Errorf("SMTP email is not configured")
	}
	port := s.cfg.SMTPPort
	if port == "" {
		port = "587"
	}
	address := s.cfg.SMTPHost + ":" + port
	headers := "From: " + s.cfg.SMTPFrom + "\r\n" + "To: " + to + "\r\n" + "Subject: " + message.subject + "\r\n" + "MIME-Version: 1.0\r\n" + "Content-Type: text/html; charset=UTF-8\r\n\r\n"

	if port == "587" {
		conn, err := smtp.Dial(address)
		if err != nil {
			return err
		}
		defer conn.Close()
		if err := conn.StartTLS(&tls.Config{ServerName: s.cfg.SMTPHost}); err != nil {
			return err
		}
		if err := conn.Auth(smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)); err != nil {
			return err
		}
		if err := conn.Mail(s.cfg.SMTPFrom); err != nil {
			return err
		}
		if err := conn.Rcpt(to); err != nil {
			return err
		}
		writer, err := conn.Data()
		if err != nil {
			return err
		}
		if _, err = writer.Write([]byte(headers + message.body)); err != nil {
			_ = writer.Close()
			return err
		}
		return writer.Close()
	}
	return smtp.SendMail(address, smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost), s.cfg.SMTPFrom, []string{to}, []byte(headers+message.body))
}
