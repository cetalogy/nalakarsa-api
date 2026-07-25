package handler

import (
	"net/http"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	res, err := h.authService.Register(req)
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

	res, err := h.authService.Login(req)
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

	res, err := h.authService.RefreshToken(req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Token refreshed successfully", res, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Logout successful", nil, nil)
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
