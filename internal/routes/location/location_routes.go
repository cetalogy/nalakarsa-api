package locationroutes

import (
	locationhandler "nalakarsa/internal/handler/location"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, h *locationhandler.LocationHandler) {
	public.GET("/locations/search", h.Search)
	public.GET("/locations/provinces", h.SearchProvinces)
	public.GET("/locations/cities", h.SearchCities)
}
