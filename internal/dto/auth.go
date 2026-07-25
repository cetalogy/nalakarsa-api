package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	Role           string `json:"role" binding:"required,oneof=Akademisi Praktisi Profesional"`
	NamaLengkap    string `json:"nama_lengkap" binding:"required"`
	GelarDepan     string `json:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi" binding:"required"`
	Lokasi         string `json:"lokasi" binding:"required"`
	BidangKeahlian string `json:"bidang_keahlian" binding:"required"`
	Industry       string `json:"industry"`
	Bio            string `json:"bio"`
	Mission        string `json:"mission"`
	AvatarURL      string `json:"avatar_url"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginData struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
