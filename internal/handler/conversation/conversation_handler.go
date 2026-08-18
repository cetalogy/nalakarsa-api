package conversationhandler

import (
	"net/http"
	"strconv"
	"strings"

	"nalakarsa/internal/dto"
	conversationservice "nalakarsa/internal/service/conversation"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConversationHandler struct {
	convService conversationservice.ConversationService
}

func NewConversationHandler(convService conversationservice.ConversationService) *ConversationHandler {
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
		utils.ValidationErrorResponse(c, err)
		return
	}

	conv, err := h.convService.GetOrCreateDirect(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Conversation retrieved successfully", conv, nil)
}

func (h *ConversationHandler) StartChat(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req dto.StartChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	conv, err := h.convService.StartChat(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Chat started successfully", conv, nil)
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
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Map Text to Body if Text is provided (for FE compatibility)
	if req.Text != "" {
		req.Body = req.Text
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

func (h *ConversationHandler) ListGroupChats(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	groupChats, err := h.convService.ListGroupChats(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Collaboration group chats retrieved successfully", groupChats, nil)
}

func (h *ConversationHandler) ListGroupMessages(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	groupIDStr := c.Param("groupId")
	if groupIDStr == "" {
		groupIDStr = c.Param("id")
	}

	groupIDStr = strings.TrimPrefix(groupIDStr, "group_topic_")
	groupIDStr = strings.TrimPrefix(groupIDStr, "group_project_")
	groupIDStr = strings.TrimPrefix(groupIDStr, "group_")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid group chat ID format: expected a valid UUID")
		return
	}

	page, limit := utils.ParsePaginationRequest(c)

	messages, total, err := h.convService.ListGroupMessages(userID, groupID, page, limit)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "group chat not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to view messages in this group chat" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Group messages retrieved successfully", messages, &dto.PaginationResponse{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		Limit:       limit,
	})
}

func (h *ConversationHandler) SendGroupMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	groupIDStr := c.Param("groupId")
	if groupIDStr == "" {
		groupIDStr = c.Param("id")
	}

	groupIDStr = strings.TrimPrefix(groupIDStr, "group_topic_")
	groupIDStr = strings.TrimPrefix(groupIDStr, "group_project_")
	groupIDStr = strings.TrimPrefix(groupIDStr, "group_")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid group chat ID format: expected a valid UUID")
		return
	}

	var req dto.SendGroupMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	msg, err := h.convService.SendGroupMessage(userID, groupID, req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "group chat not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to send messages to this group chat" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Group message sent successfully", msg, nil)
}

func (h *ConversationHandler) DeleteMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	msgIDStr := c.Param("id")
	if msgIDStr == "" {
		msgIDStr = c.Param("messageId")
	}

	msgID, err := uuid.Parse(msgIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid message ID format: expected a valid UUID")
		return
	}

	if err := h.convService.DeleteMessage(userID, msgID); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "message not found" {
			statusCode = http.StatusNotFound
		} else if strings.HasPrefix(err.Error(), "unauthorized") {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Message deleted successfully", nil, nil)
}

