package handler

import (
	"net/http"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiscussionHandler struct {
	discService service.DiscussionService
}

func NewDiscussionHandler(discService service.DiscussionService) *DiscussionHandler {
	return &DiscussionHandler{discService: discService}
}

// getCurrentUserID tries to extract user_id from context (may be nil for public endpoints)
func getCurrentUserID(c *gin.Context) *uuid.UUID {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return nil
	}
	uid := userIDInterface.(uuid.UUID)
	return &uid
}

func (h *DiscussionHandler) Create(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	var req dto.CreateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	id, err := h.discService.Create(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Discussion topic created successfully", gin.H{"id": id}, nil)
}

func (h *DiscussionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid discussion ID format")
		return
	}

	currentUserID := getCurrentUserID(c)

	disc, err := h.discService.GetByID(id, currentUserID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Discussion topic retrieved successfully", disc, nil)
}

func (h *DiscussionHandler) List(c *gin.Context) {
	search := c.Query("q")
	if search == "" {
		search = c.Query("search")
	}
	category := c.Query("category")
	role := c.Query("role")
	sort := c.Query("sort")
	page, limit := utils.ParsePaginationRequest(c)

	currentUserID := getCurrentUserID(c)

	discs, total, err := h.discService.List(search, category, role, sort, page, limit, currentUserID)
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

	utils.JSONResponse(c, http.StatusOK, "Discussions list retrieved successfully", discs, pagination)
}

func (h *DiscussionHandler) Update(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid discussion ID format")
		return
	}

	var req dto.UpdateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	if err := h.discService.Update(userID, id, req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "discussion not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to update this discussion" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Discussion topic updated successfully", nil, nil)
}

func (h *DiscussionHandler) Delete(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid discussion ID format")
		return
	}

	if err := h.discService.Delete(userID, id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "discussion not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to delete this discussion" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Discussion topic deleted successfully", nil, nil)
}

func (h *DiscussionHandler) AddReply(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	discussionID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid discussion ID format")
		return
	}

	var req dto.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	reply, err := h.discService.AddReply(userID, discussionID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "discussion not found" {
			statusCode = http.StatusNotFound
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Reply added successfully", reply, nil)
}

func (h *DiscussionHandler) DeleteReply(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	replyIDStr := c.Param("reply_id")
	replyID, err := uuid.Parse(replyIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid reply ID format")
		return
	}

	if err := h.discService.DeleteReply(userID, replyID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reply not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to delete this reply" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Reply deleted successfully", nil, nil)
}

func (h *DiscussionHandler) Vote(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	discussionID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid discussion ID format")
		return
	}

	if err := h.discService.Vote(userID, discussionID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "discussion not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "already upvoted this discussion" {
			statusCode = http.StatusConflict
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Upvoted successfully", nil, nil)
}

func (h *DiscussionHandler) Unvote(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	discussionID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid discussion ID format")
		return
	}

	if err := h.discService.Unvote(userID, discussionID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Upvote removed successfully", nil, nil)
}
