package handler

import (
	"net/http"
	"strconv"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConversationHandler struct {
	convService service.ConversationService
}

func NewConversationHandler(convService service.ConversationService) *ConversationHandler {
	return &ConversationHandler{convService: convService}
}

func (h *ConversationHandler) ListConversations(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, limit := utils.ParsePaginationRequest(c)

	convs, total, err := h.convService.ListConversations(userID, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Conversations retrieved successfully", convs, &dto.PaginationResponse{
		CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit,
	})
}

func (h *ConversationHandler) GetOrCreateDirect(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req dto.CreateDirectConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	conv, err := h.convService.GetOrCreateDirect(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Conversation retrieved successfully", conv, nil)
}

func (h *ConversationHandler) ListMessages(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid conversation ID format")
		return
	}

	limit := 20
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil && l > 0 {
		limit = l
	}
	cursor := c.Query("cursor")

	messages, hasMore, err := h.convService.ListMessages(userID, convID, limit, cursor)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Messages retrieved successfully", gin.H{
		"messages": messages,
		"has_more": hasMore,
	}, nil)
}

func (h *ConversationHandler) SendMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid conversation ID format")
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	msg, err := h.convService.SendMessage(userID, convID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Message sent successfully", msg, nil)
}

func (h *ConversationHandler) MarkRead(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid conversation ID format")
		return
	}

	if err := h.convService.MarkRead(userID, convID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Messages marked as read", nil, nil)
}
