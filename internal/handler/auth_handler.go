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
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	res, err := h.authService.Register(req)
	if err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Registration successful", res, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		utils.ErrorJSONResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Login successful", res, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	res, err := h.authService.RefreshToken(req)
	if err != nil {
		utils.ErrorJSONResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Token refreshed successfully", res, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		utils.ErrorJSONResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Logout successful", nil, nil)
}
