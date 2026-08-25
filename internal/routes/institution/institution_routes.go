package institutionroutes

import (
	institutionhandler "nalakarsa/internal/handler/institution"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, h *institutionhandler.InstitutionHandler) {
	public.GET("/institutions/search", h.Search)
}
