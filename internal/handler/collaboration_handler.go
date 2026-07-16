package handler

import (
	"net/http"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollaborationHandler struct {
	collabService service.CollaborationService
}

func NewCollaborationHandler(collabService service.CollaborationService) *CollaborationHandler {
	return &CollaborationHandler{collabService: collabService}
}

func (h *CollaborationHandler) Create(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	var req dto.CreateCollaborationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	id, err := h.collabService.Create(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Collaboration proposal created successfully", gin.H{"id": id}, nil)
}

func (h *CollaborationHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid collaboration proposal ID format")
		return
	}

	collab, err := h.collabService.GetByID(id)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Collaboration proposal retrieved successfully", collab, nil)
}

func (h *CollaborationHandler) List(c *gin.Context) {
	search := c.Query("search")
	roleRequired := c.Query("role_required")
	status := c.Query("status")
	page, limit := utils.ParsePaginationRequest(c)

	collabs, total, err := h.collabService.List(search, roleRequired, status, page, limit)
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

	utils.JSONResponse(c, http.StatusOK, "Collaboration proposals list retrieved successfully", collabs, pagination)
}

func (h *CollaborationHandler) Update(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid collaboration proposal ID format")
		return
	}

	var req dto.UpdateCollaborationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	if err := h.collabService.Update(userID, id, req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "collaboration proposal not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to update this collaboration proposal" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Collaboration proposal updated successfully", nil, nil)
}

func (h *CollaborationHandler) Delete(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid collaboration proposal ID format")
		return
	}

	if err := h.collabService.Delete(userID, id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "collaboration proposal not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to delete this collaboration proposal" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Collaboration proposal deleted successfully", nil, nil)
}

func (h *CollaborationHandler) Apply(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid collaboration proposal ID format")
		return
	}

	var req dto.ApplyCollaborationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	appID, err := h.collabService.Apply(userID, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "collaboration proposal not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot apply to your own collaboration proposal" ||
			err.Error() == "collaboration proposal is no longer open" ||
			err.Error() == "your role does not match the required role for this collaboration" ||
			err.Error() == "you have already applied to this collaboration proposal" {
			statusCode = http.StatusBadRequest
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Collaboration application submitted successfully", gin.H{"id": appID}, nil)
}

func (h *CollaborationHandler) ListApplications(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid collaboration proposal ID format")
		return
	}

	apps, err := h.collabService.ListApplications(userID, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "collaboration proposal not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to view applications for this proposal" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Applications list retrieved successfully", apps, nil)
}

func (h *CollaborationHandler) UpdateApplicationStatus(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid collaboration proposal ID format")
		return
	}

	appIDStr := c.Param("app_id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid application ID format")
		return
	}

	var req dto.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	if err := h.collabService.UpdateApplicationStatus(userID, id, appID, req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "collaboration proposal not found" || err.Error() == "application not found for this proposal" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to manage applications for this proposal" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Application status updated successfully", nil, nil)
}
