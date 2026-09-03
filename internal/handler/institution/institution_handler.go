package institutionhandler

import (
	"net/http"
	"strconv"
	"strings"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/service/institution"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type InstitutionHandler struct {
	institutionService institutionservice.InstitutionService
}

func (h *InstitutionHandler) Create(c *gin.Context) {
	var req dto.CreateReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "name is required")
		return
	}
	item, err := h.institutionService.CreateInstitution(req.Name)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Institution created successfully", item, nil)
}

func NewInstitutionHandler(institutionService institutionservice.InstitutionService) *InstitutionHandler {
	return &InstitutionHandler{
		institutionService: institutionService,
	}
}

func (h *InstitutionHandler) Search(c *gin.Context) {
	search := strings.TrimSpace(c.Query("q"))

	page := 1
	if pageQuery := c.Query("page"); pageQuery != "" {
		parsedPage, err := strconv.Atoi(pageQuery)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if limit > 20 {
		limit = 20
	}

	institutions, total, err := h.institutionService.SearchInstitutions(search, page, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	utils.JSONResponse(c, http.StatusOK, "Institution suggestions retrieved successfully", institutions, &dto.PaginationResponse{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		Limit:       limit,
	})
}
