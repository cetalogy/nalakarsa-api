package homeroutes

import (
	homepagehandler "nalakarsa/internal/handler/homepage"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, h *homepagehandler.HomepageHandler) {
	public.GET("/homepage/hero", h.GetHero)
	public.GET("/homepage/sections", h.GetSections)
	public.GET("/homepage/testimonials", h.GetTestimonials)
}
