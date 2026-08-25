package userhandler

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	usercommon "nalakarsa/internal/common/user"
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
	var userID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uuid.UUID)
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				userID = claims.UserID
			} else {
				utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
				return
			}
		} else {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Profile retrieved successfully", profile, nil)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var userID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uuid.UUID)
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				userID = claims.UserID
			} else {
				utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
				return
			}
		} else {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

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
	if idStr == "" || idStr == "me" {
		h.GetProfile(c)
		return
	}

	targetUser, err := h.userService.ResolveUser(idStr)
	if err != nil || targetUser == nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, "User not found")
		return
	}

	profile, err := h.userService.GetPublicProfile(targetUser.ID)
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
	var userID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uuid.UUID)
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				userID = claims.UserID
			} else {
				utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
				return
			}
		} else {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

	projects, err := h.userService.GetMyProjects(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "My projects retrieved successfully", projects, nil)
}

func (h *UserHandler) GetMyStats(c *gin.Context) {
	var userID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uuid.UUID)
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				userID = claims.UserID
			} else {
				utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
				return
			}
		} else {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

	stats, err := h.userService.GetMyStats(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "My stats retrieved successfully", stats, nil)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, usercommon.MaxAvatarRequestSize)
	var userID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uuid.UUID)
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				userID = claims.UserID
			} else {
				utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
				return
			}
		} else {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

	// Get file from form (flexible key: avatar, file, image, photo)
	file, err := c.FormFile("avatar")
	if err != nil {
		file, err = c.FormFile("file")
	}
	if err != nil {
		file, err = c.FormFile("image")
	}
	if err != nil {
		file, err = c.FormFile("photo")
	}
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "File 'avatar' (or 'file' / 'image') is required")
		return
	}

	log.Printf("[AVATAR UPLOAD] User %s uploading file '%s' (%d bytes, Content-Type: %s)", userID, file.Filename, file.Size, file.Header.Get("Content-Type"))

	// Validate file size (max 5MB)
	if file.Size > usercommon.MaxAvatarSize {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "File size exceeds 5MB limit")
		return
	}

	// Validate file type (strictly JPG, JPEG, PNG)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedExts[ext] || (contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/png") {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid file type. Only JPG, JPEG, and PNG are allowed")
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Unable to read uploaded image")
		return
	}
	defer openedFile.Close()
	header := make([]byte, 512)
	readCount, err := openedFile.Read(header)
	if err != nil && err != io.EOF {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Unable to inspect uploaded image")
		return
	}
	detectedType := http.DetectContentType(header[:readCount])
	if detectedType != "image/jpeg" && detectedType != "image/png" {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Uploaded file is not a valid JPG or PNG image")
		return
	}

	// Upload to Supabase Storage
	secureURL, err := utils.UploadAvatarToSupabase(userID, file, h.cfg)
	if err != nil {
		log.Printf("[AVATAR UPLOAD FAILED] User %s: %v", userID, err)
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Failed to upload avatar to cloud storage: "+err.Error())
		return
	}

	// Update in database
	if err := h.userService.UpdateAvatar(userID, secureURL); err != nil {
		log.Printf("[AVATAR DB UPDATE FAILED] User %s: %v", userID, err)
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, "Failed to update profile avatar: "+err.Error())
		return
	}

	if err := utils.CleanupOtherAvatarVariants(userID, secureURL, h.cfg); err != nil {
		log.Printf("failed to cleanup old avatar variants for user %s: %v", userID, err)
	}

	log.Printf("[AVATAR UPLOAD SUCCESS] User %s avatar set to %s", userID, secureURL)

	utils.JSONResponse(c, http.StatusOK, "Avatar uploaded successfully", gin.H{
		"avatar_url": secureURL,
		"avatarUrl":  secureURL,
	}, nil)
}

func (h *UserHandler) ToggleFollow(c *gin.Context) {
	var currentUserID uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		currentUserID = userIDInterface.(uuid.UUID)
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				currentUserID = claims.UserID
			} else {
				utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
				return
			}
		} else {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}

	targetIDStr := c.Param("id")
	targetUser, err := h.userService.ResolveUser(targetIDStr)
	if err != nil || targetUser == nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, "Target user not found")
		return
	}

	res, err := h.userService.ToggleFollow(currentUserID, targetUser.ID)
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
	var targetUserID uuid.UUID

	if targetIDStr == "" || targetIDStr == "me" {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			// Fallback: parse Bearer token from header if hit via public group
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
				if err == nil {
					targetUserID = claims.UserID
					exists = true
				}
			}
		} else {
			targetUserID = userIDInterface.(uuid.UUID)
		}

		if !exists {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	} else {
		targetUser, err := h.userService.ResolveUser(targetIDStr)
		if err != nil || targetUser == nil {
			utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, "User not found")
			return
		}
		targetUserID = targetUser.ID
	}

	var currentUserID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		uid := userIDInterface.(uuid.UUID)
		currentUserID = &uid
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				currentUserID = &claims.UserID
			}
		}
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
	var targetUserID uuid.UUID

	if targetIDStr == "" || targetIDStr == "me" {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			// Fallback: parse Bearer token from header if hit via public group
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
				if err == nil {
					targetUserID = claims.UserID
					exists = true
				}
			}
		} else {
			targetUserID = userIDInterface.(uuid.UUID)
		}

		if !exists {
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
	} else {
		targetUser, err := h.userService.ResolveUser(targetIDStr)
		if err != nil || targetUser == nil {
			utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, "User not found")
			return
		}
		targetUserID = targetUser.ID
	}

	var currentUserID *uuid.UUID
	if userIDInterface, exists := c.Get("user_id"); exists {
		uid := userIDInterface.(uuid.UUID)
		currentUserID = &uid
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenStr, h.cfg.JWTSecret)
			if err == nil {
				currentUserID = &claims.UserID
			}
		}
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
