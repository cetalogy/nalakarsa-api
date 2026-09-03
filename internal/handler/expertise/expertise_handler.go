package expertisehandler

import (
	"net/http"
	"strconv"
	"strings"

	"nalakarsa/internal/dto"
	expertiseservice "nalakarsa/internal/service/expertise"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type ExpertiseHandler struct {
	service expertiseservice.ExpertiseService
}

func (h *ExpertiseHandler) Create(c *gin.Context) {
	var req dto.CreateReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "name is required")
		return
	}
	item, err := h.service.Create(req.Name)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Expertise created successfully", item, nil)
}

func NewExpertiseHandler(service expertiseservice.ExpertiseService) *ExpertiseHandler {
	return &ExpertiseHandler{service: service}
}

func (h *ExpertiseHandler) Search(c *gin.Context) {
	limit := 10
	if value := c.Query("limit"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	result, err := h.service.Search(strings.TrimSpace(c.Query("q")), limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Expertise suggestions retrieved successfully", result, nil)
}
