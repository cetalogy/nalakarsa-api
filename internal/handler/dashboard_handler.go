package handler

import (
	"net/http"

	"nalakarsa/internal/service"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashService service.DashboardService
}

func NewDashboardHandler(dashService service.DashboardService) *DashboardHandler {
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
