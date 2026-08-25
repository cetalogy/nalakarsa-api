package expertiseroutes

import (
	expertisehandler "nalakarsa/internal/handler/expertise"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, handler *expertisehandler.ExpertiseHandler) {
	public.GET("/expertise/search", handler.Search)
}
