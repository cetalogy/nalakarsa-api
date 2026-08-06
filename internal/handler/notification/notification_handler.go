package notificationhandler

import (
	"net/http"

	"nalakarsa/internal/dto"
	notificationservice "nalakarsa/internal/service/notification"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notifService notificationservice.NotificationService
}

func NewNotificationHandler(notifService notificationservice.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, limit := utils.ParsePaginationRequest(c)

	notifs, total, err := h.notifService.ListNotifications(userID, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Notifications retrieved successfully", notifs, &dto.PaginationResponse{
		CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit,
	})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	notifID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid notification ID format")
		return
	}

	if err := h.notifService.MarkRead(userID, notifID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Notification marked as read", nil, nil)
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.notifService.MarkAllRead(userID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "All notifications marked as read", nil, nil)
}
