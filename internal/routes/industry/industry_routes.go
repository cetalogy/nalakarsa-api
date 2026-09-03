package industryroutes

import (
	"github.com/gin-gonic/gin"
	industryhandler "nalakarsa/internal/handler/industry"
)

func RegisterRoutes(public *gin.RouterGroup, handler *industryhandler.IndustryHandler) {
	public.GET("/industries/search", handler.Search)
	public.POST("/industries", handler.Create)
}
