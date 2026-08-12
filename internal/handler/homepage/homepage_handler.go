package homepagehandler

import (
	"net/http"

	homepageService "nalakarsa/internal/service/homepage"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

type HomepageHandler struct {
	homepageService homepageService.HomepageService
}

func NewHomepageHandler(homepageService homepageService.HomepageService) *HomepageHandler {
	return &HomepageHandler{
		homepageService: homepageService,
	}
}

func (h *HomepageHandler) GetHero(c *gin.Context) {
	hero, err := h.homepageService.GetHero()
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Homepage hero retrieved successfully", hero, nil)
}

func (h *HomepageHandler) GetSections(c *gin.Context) {
	sections, err := h.homepageService.GetSections()
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Homepage sections retrieved successfully", sections, nil)
}

func (h *HomepageHandler) GetTestimonials(c *gin.Context) {
	testimonials, err := h.homepageService.GetTestimonials()
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Homepage testimonials retrieved successfully", testimonials, nil)
}
