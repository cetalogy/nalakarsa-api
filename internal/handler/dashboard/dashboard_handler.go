package dashboardhandler

import (
	"net/http"

	"nalakarsa/internal/dto"
	dashboardservice "nalakarsa/internal/service/dashboard"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashService dashboardservice.DashboardService
}

func NewDashboardHandler(dashService dashboardservice.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashService: dashService}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	dashboard, err := h.dashService.GetDashboard(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Dashboard retrieved successfully", dashboard, nil)
}

func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	summary, err := h.dashService.GetDashboardSummary(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Dashboard summary retrieved successfully", summary, nil)
}

func (h *DashboardHandler) GetDashboardOngoingProjects(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	projects, err := h.dashService.GetDashboardOngoingProjects(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Dashboard ongoing projects retrieved successfully", projects, nil)
}

func (h *DashboardHandler) GetDashboardNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	page, limit := utils.ParsePaginationRequest(c)

	notifs, total, err := h.dashService.GetDashboardNotifications(userID, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.JSONResponse(c, http.StatusOK, "Dashboard notifications retrieved successfully", notifs, &dto.PaginationResponse{
		CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit,
	})
}
