package handler

import (
	"net/http"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projService service.ProjectService
}

func NewProjectHandler(projService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projService: projService}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	id, err := h.projService.Create(userID, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Project created successfully", gin.H{"id": id}, nil)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	project, err := h.projService.GetByID(id)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Project retrieved successfully", project, nil)
}

func (h *ProjectHandler) List(c *gin.Context) {
	search := c.Query("q")
	if search == "" {
		search = c.Query("search")
	}
	status := c.Query("status")
	category := c.Query("category")
	page, limit := utils.ParsePaginationRequest(c)

	projects, total, err := h.projService.List(search, status, category, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Projects list retrieved successfully", projects, &dto.PaginationResponse{
		CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit,
	})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	if err := h.projService.Update(userID, id, req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "project not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to update this project" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Project updated successfully", nil, nil)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	if err := h.projService.Delete(userID, id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "project not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to delete this project" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Project deleted successfully", nil, nil)
}

func (h *ProjectHandler) Apply(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	var req dto.ApplyProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	appID, err := h.projService.Apply(userID, id, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Application submitted successfully", gin.H{"id": appID}, nil)
}

func (h *ProjectHandler) ListApplications(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	apps, err := h.projService.ListApplications(userID, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "project not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to view applications for this project" {
			statusCode = http.StatusForbidden
		}
		utils.ErrorJSONResponseWithMessage(c, statusCode, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Applications retrieved successfully", apps, nil)
}

func (h *ProjectHandler) UpdateApplicationStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}
	appID, err := uuid.Parse(c.Param("applicationId"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid application ID format")
		return
	}

	var req dto.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	if err := h.projService.UpdateApplicationStatus(userID, id, appID, req); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Application status updated successfully", nil, nil)
}

func (h *ProjectHandler) ListMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	members, err := h.projService.ListMembers(id)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Members retrieved successfully", members, nil)
}

func (h *ProjectHandler) CreateMilestone(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	var req dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	msID, err := h.projService.CreateMilestone(userID, id, req)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "Milestone created successfully", gin.H{"id": msID}, nil)
}

func (h *ProjectHandler) UpdateMilestone(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}
	msID, err := uuid.Parse(c.Param("milestoneId"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "Invalid milestone ID format")
		return
	}

	var req dto.UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", nil)
		return
	}

	if err := h.projService.UpdateMilestone(userID, id, msID, req); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Milestone updated successfully", nil, nil)
}
