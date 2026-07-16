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

func (h *DiscussionHandler) Create(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	var req dto.CreateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
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

	disc, err := h.discService.GetByID(id)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Discussion topic retrieved successfully", disc, nil)
}

func (h *DiscussionHandler) List(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	role := c.Query("role")
	sort := c.Query("sort")
	page, limit := utils.ParsePaginationRequest(c)

	discs, total, err := h.discService.List(search, category, role, sort, page, limit)
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
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
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

func (h *DiscussionHandler) AddComment(c *gin.Context) {
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

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "Validation failed", []string{err.Error()})
		return
	}

	comment, err := h.discService.AddComment(userID, discussionID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "discussion not found" {
			statusCode = http.StatusNotFound
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Comment added successfully", comment, nil)
}

func (h *DiscussionHandler) DeleteComment(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDInterface.(uuid.UUID)

	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid comment ID format")
		return
	}

	if err := h.discService.DeleteComment(userID, commentID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to delete this comment" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Comment deleted successfully", nil, nil)
}
