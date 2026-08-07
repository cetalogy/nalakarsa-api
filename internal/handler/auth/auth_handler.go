package authhandler

import (
	"net/http"
	"log"
	"strings"

	"nalakarsa/internal/dto"
	authservice "nalakarsa/internal/service/auth"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService authservice.AuthService
}

func NewAuthHandler(authService authservice.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	res, err := h.authService.Register(req, h.buildRequestContext(c))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Registration successful", res, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	res, err := h.authService.Login(req, h.buildRequestContext(c))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Login successful", res, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	res, err := h.authService.RefreshToken(req, h.buildRequestContext(c))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Token refreshed successfully", res, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	// FE may not always send refresh_token for logout
	_ = c.ShouldBindJSON(&req)

	// Audit log always exists for logout events
	userID, _ := c.Get("user_id")
	if req.RefreshToken != "" {
		log.Printf("[auth][logout] user_id=%v ip=%s action=logout_refresh_revoke requested=true", userID, c.ClientIP())
		if err := h.authService.Logout(req.RefreshToken); err != nil {
			log.Printf("[auth][logout] user_id=%v ip=%s action=logout_refresh_revoke result=failed error=%v", userID, c.ClientIP(), err)
			// keep response success to avoid token-existence leakage
		} else {
			log.Printf("[auth][logout] user_id=%v ip=%s action=logout_refresh_revoke result=success", userID, c.ClientIP())
		}
		utils.JSONResponse(c, http.StatusOK, "Logout successful", nil, nil)
		return
	}

	log.Printf("[auth][logout] user_id=%v ip=%s action=logout request_no_token", userID, c.ClientIP())
	utils.JSONResponse(c, http.StatusOK, "Logout successful", nil, nil)
}

func (h *AuthHandler) buildRequestContext(c *gin.Context) *dto.AuthRequestContext {
	deviceInfo := strings.TrimSpace(c.GetHeader("X-Device-Info"))
	if deviceInfo == "" {
		deviceInfo = strings.TrimSpace(c.GetHeader("User-Agent"))
	}

	return &dto.AuthRequestContext{
		DeviceInfo: deviceInfo,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	// TODO: Implement email sending for password reset
	// For now, return success to prevent email enumeration
	utils.JSONResponse(c, http.StatusOK, "If the email exists, a reset link has been sent", nil, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	// TODO: Implement token validation and password reset
	utils.ErrorJSONResponseWithMessage(c, http.StatusNotImplemented, "Password reset not yet implemented")
}
