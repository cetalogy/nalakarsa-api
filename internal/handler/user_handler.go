package handler

import (
	"net/http"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
	cfg         *config.Config
}

func NewUserHandler(userService service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Profile retrieved successfully", profile, nil)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	if err := h.userService.UpdateProfile(userID, req); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Profile updated successfully", nil, nil)
}

func (h *UserHandler) GetPublicProfile(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	profile, err := h.userService.GetPublicProfile(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Public profile retrieved successfully", profile, nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	search := c.Query("q")
	if search == "" {
		search = c.Query("search")
	}
	role := c.Query("role")
	page, limit := utils.ParsePaginationRequest(c)

	users, total, err := h.userService.ListUsers(search, role, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := &dto.PaginationResponse{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		Limit:       limit,
	}

	utils.JSONResponse(c, http.StatusOK, "Users list retrieved successfully", users, pagination)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "File 'avatar' is required")
		return
	}

	// Validate file size (max 2MB)
	const maxFileSize = 2 * 1024 * 1024 // 2MB
	if file.Size > maxFileSize {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "File size exceeds 2MB limit")
		return
	}

	// Validate file type (must be image)
	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, GIF, and WEBP are allowed")
		return
	}

	// Upload to Firebase
	secureURL, err := utils.UploadAvatarToFirebase(file, h.cfg)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Failed to upload avatar to cloud storage: "+err.Error())
		return
	}

	// Update in database
	if err := h.userService.UpdateAvatar(userID, secureURL); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Failed to update profile avatar: "+err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Avatar uploaded successfully", gin.H{
		"avatar_url": secureURL,
	}, nil)
}
