package organizationhandler

import (
	"net/http"
	"strconv"
	"strings"

	"nalakarsa/internal/dto"
	organizationservice "nalakarsa/internal/service/organization"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service organizationservice.OrganizationService
}

func NewOrganizationHandler(service organizationservice.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) Search(c *gin.Context) {
	page, limit := 1, 10
	if value, err := strconv.Atoi(c.Query("page")); err == nil && value > 0 {
		page = value
	}
	if value, err := strconv.Atoi(c.Query("limit")); err == nil && value > 0 && value <= 20 {
		limit = value
	}
	items, total, err := h.service.Search(strings.TrimSpace(c.Query("q")), page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	utils.JSONResponse(c, http.StatusOK, "Organizations retrieved successfully", items, &dto.PaginationResponse{CurrentPage: page, TotalPages: totalPages, TotalItems: total, Limit: limit})
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	item, err := h.service.Create(req.Name)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Organization created successfully", item, nil)
}
