package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email          string `json:"email" form:"email" binding:"required,email"`
	Password       string `json:"password" form:"password" binding:"required,min=8"`
	Role           string `json:"role" form:"role" binding:"required,oneof=Akademisi Praktisi Profesional"`
	NamaLengkap    string `json:"nama_lengkap" form:"nama_lengkap" binding:"required"`
	GelarDepan     string `json:"gelar_depan" form:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang" form:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi" form:"afiliasi" binding:"required"`
	Lokasi         string `json:"lokasi" form:"lokasi" binding:"required"`
	BidangKeahlian string `json:"bidang_keahlian" form:"bidang_keahlian" binding:"required"`
	Industry       string `json:"industry" form:"industry"`
	Bio            string `json:"bio" form:"bio"`
	Mission        string `json:"mission" form:"mission"`
	AvatarURL      string `json:"avatar_url" form:"avatar_url"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id" form:"id"`
	Email string    `json:"email" form:"email"`
	Role  string    `json:"role" form:"role"`
}

type RegisterResponse struct {
	ID    uuid.UUID `json:"id" form:"id"`
	Email string    `json:"email" form:"email"`
	Role  string    `json:"role" form:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginData struct {
	AccessToken  string       `json:"access_token" form:"access_token"`
	RefreshToken string       `json:"refresh_token" form:"refresh_token"`
	ExpiresIn    int          `json:"expires_in" form:"expires_in"`
	User         UserResponse `json:"user" form:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type RefreshTokenData struct {
	AccessToken  string `json:"access_token" form:"access_token"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	ExpiresIn    int    `json:"expires_in" form:"expires_in"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" form:"token" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required,min=8"`
}
