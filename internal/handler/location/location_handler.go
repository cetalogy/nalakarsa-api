package locationhandler

import (
	"net/http"
	"strconv"
	"strings"

	"nalakarsa/internal/service/location"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationService locationservice.LocationService
}

func NewLocationHandler(locationService locationservice.LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

func (h *LocationHandler) Search(c *gin.Context) {
	h.searchWithType(c, strings.TrimSpace(c.Query("type")))
}

func (h *LocationHandler) SearchProvinces(c *gin.Context) {
	h.searchWithType(c, "province")
}

func (h *LocationHandler) SearchCities(c *gin.Context) {
	h.searchWithType(c, "city")
}

func (h *LocationHandler) searchWithType(c *gin.Context, forcedType string) {
	search := strings.TrimSpace(c.Query("q"))
	locationType := strings.ToLower(strings.TrimSpace(forcedType))
	if locationType == "" {
		locationType = strings.ToLower(strings.TrimSpace(c.Query("type")))
	}
	provinceID := strings.TrimSpace(c.Query("provinceId"))

	limit := 10
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	locations, err := h.locationService.SearchLocations(search, locationType, provinceID, limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Locations retrieved successfully", locations, nil)
}
