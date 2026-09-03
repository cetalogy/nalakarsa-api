package industryroutes

import (
	industryhandler "nalakarsa/internal/handler/industry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, handler *industryhandler.IndustryHandler) {
	public.GET("/industries/search", handler.Search)
}
