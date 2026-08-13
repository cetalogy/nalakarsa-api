package connectionhandler

import (
	"net/http"

	"nalakarsa/internal/dto"
	connectionservice "nalakarsa/internal/service/connection"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConnectionHandler struct {
	connService connectionservice.ConnectionService
}

func NewConnectionHandler(connService connectionservice.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{connService: connService}
}

func (h *ConnectionHandler) ListConnections(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, limit := utils.ParsePaginationRequest(c)

	conns, total, err := h.connService.ListConnections(userID, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Connections retrieved successfully", conns, &dto.PaginationResponse{
		CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit,
	})
}

func (h *ConnectionHandler) ListRequests(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	requestType := c.DefaultQuery("type", "incoming")
	page, limit := utils.ParsePaginationRequest(c)

	reqs, total, err := h.connService.ListRequests(userID, requestType, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Connection requests retrieved successfully", reqs, &dto.PaginationResponse{
		CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit,
	})
}

func (h *ConnectionHandler) SendRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req dto.SendConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	id, err := h.connService.SendRequest(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Connection request sent successfully", gin.H{"id": id}, nil)
}

func (h *ConnectionHandler) AcceptRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid request ID format")
		return
	}

	if err := h.connService.AcceptRequest(userID, requestID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Connection request accepted", nil, nil)
}

func (h *ConnectionHandler) RejectRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid request ID format")
		return
	}

	if err := h.connService.RejectRequest(userID, requestID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Connection request rejected", nil, nil)
}

func (h *ConnectionHandler) CancelRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid request ID format")
		return
	}

	if err := h.connService.CancelRequest(userID, requestID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Connection request cancelled", nil, nil)
}

func (h *ConnectionHandler) RemoveConnection(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	if err := h.connService.RemoveConnection(userID, targetUserID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Connection removed successfully", nil, nil)
}

func (h *ConnectionHandler) GetSuggestions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	suggestions, err := h.connService.GetSuggestions(userID, 10)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Suggestions retrieved successfully", suggestions, nil)
}
