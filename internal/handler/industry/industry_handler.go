package industryhandler

import (
	"net/http"
	"strconv"
	"strings"

	industryservice "nalakarsa/internal/service/industry"
	"nalakarsa/internal/dto"
	"nalakarsa/internal/utils"
	"github.com/gin-gonic/gin"
)

type IndustryHandler struct{ service industryservice.IndustryService }

func NewIndustryHandler(service industryservice.IndustryService) *IndustryHandler {
	return &IndustryHandler{service: service}
}

func (h *IndustryHandler) Create(c *gin.Context) {
	var req dto.CreateReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "name is required"); return }
	item, err := h.service.Create(req.Name)
	if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error()); return }
	utils.JSONResponse(c, http.StatusCreated, "Industry created successfully", item, nil)
}

func (h *IndustryHandler) Search(c *gin.Context) {
	limit := 10
	if value := c.Query("limit"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 { limit = parsed }
	}
	result, err := h.service.Search(strings.TrimSpace(c.Query("q")), limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Industry suggestions retrieved successfully", result, nil)
}
