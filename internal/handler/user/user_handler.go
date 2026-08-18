package userhandler

import (
	"log"
	"net/http"

	"nalakarsa/internal/config"
	"nalakarsa/internal/dto"
	userservice "nalakarsa/internal/service/user"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService userservice.UserService
	cfg         *config.Config
}

func NewUserHandler(userService userservice.UserService, cfg *config.Config) *UserHandler {
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
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.userService.UpdateProfile(userID, req); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.JSONResponse(c, http.StatusOK, "Profile updated successfully", nil, nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Profile updated successfully", profile, nil)
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

func (h *UserHandler) GetMyProjects(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	projects, err := h.userService.GetMyProjects(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "My projects retrieved successfully", projects, nil)
}

func (h *UserHandler) GetMyStats(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	stats, err := h.userService.GetMyStats(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "My stats retrieved successfully", stats, nil)
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
	}
	if !allowedTypes[contentType] {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid file type. Only JPG, JPEG, and PNG are allowed")
		return
	}

	// Upload to Supabase Storage
	secureURL, err := utils.UploadAvatarToSupabase(userID, file, h.cfg)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Failed to upload avatar to cloud storage: "+err.Error())
		return
	}

	// Update in database
	if err := h.userService.UpdateAvatar(userID, secureURL); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Failed to update profile avatar: "+err.Error())
		return
	}

	if err := utils.CleanupOtherAvatarVariants(userID, secureURL, h.cfg); err != nil {
		log.Printf("failed to cleanup old avatar variants for user %s: %v", userID, err)
	}

	utils.JSONResponse(c, http.StatusOK, "Avatar uploaded successfully", gin.H{
		"avatar_url": secureURL,
	}, nil)
}

func (h *UserHandler) ToggleFollow(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	currentUserID := userIDInterface.(uuid.UUID)

	targetIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid target user ID format")
		return
	}

	res, err := h.userService.ToggleFollow(currentUserID, targetUserID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "target user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot follow yourself" {
			statusCode = http.StatusBadRequest
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, res.Message, res, nil)
}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var currentUserID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		uid := userIDInterface.(uuid.UUID)
		currentUserID = &uid
	}

	page, limit := utils.ParsePaginationRequest(c)

	followers, total, err := h.userService.GetFollowers(currentUserID, targetUserID, page, limit)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
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

	utils.JSONResponse(c, http.StatusOK, "Followers retrieved successfully", followers, pagination)
}

func (h *UserHandler) GetFollowing(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetUserID, err := uuid.Parse(targetIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var currentUserID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		uid := userIDInterface.(uuid.UUID)
		currentUserID = &uid
	}

	page, limit := utils.ParsePaginationRequest(c)

	following, total, err := h.userService.GetFollowing(currentUserID, targetUserID, page, limit)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
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

	utils.JSONResponse(c, http.StatusOK, "Following list retrieved successfully", following, pagination)
}

