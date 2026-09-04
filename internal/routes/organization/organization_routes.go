package organizationroutes

import (
	organizationhandler "nalakarsa/internal/handler/organization"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, handler *organizationhandler.OrganizationHandler) {
	public.GET("/organizations/search", handler.Search)
	public.POST("/organizations", handler.Create)
}
