package dto

type RegisterRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6,max=72"`
	SecurityQuestion string `json:"security_question" binding:"required,oneof=favorite_animal mother_name"`
	SecurityAnswer string `json:"security_answer" binding:"required,min=2,max=100"`
	Role        string  `json:"role" binding:"required"`
	FirstName   string  `json:"firstName" binding:"required"`
	MiddleName  *string `json:"middleName"`
	LastName    string  `json:"lastName" binding:"required"`
	FullName    string  `json:"fullName" binding:"required"`
	PrefixTitle string  `json:"prefixTitle"`
	SuffixTitle string  `json:"suffixTitle"`
	Affiliation string  `json:"affiliation" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	Expertise   string  `json:"expertise" binding:"required"`
	Industry    string  `json:"industry"`
	Bio         string  `json:"bio" binding:"max=500"`
	AvatarURL   string  `json:"avatar_url"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthData struct {
	Token        string       `json:"token"` // For FE spec
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type AuthRequestContext struct {
	DeviceInfo string
	IPAddress  string
	UserAgent  string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
	SecurityQuestion string `json:"security_question" binding:"required,oneof=favorite_animal mother_name"`
	SecurityAnswer string `json:"security_answer" binding:"required,min=2,max=100"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=72"`
}
