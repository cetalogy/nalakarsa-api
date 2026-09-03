package knowledgeroutes

import (
	knowledgehandler "nalakarsa/internal/handler/knowledge"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, handler *knowledgehandler.KnowledgeHandler) {
	public.GET("/knowledge/search", handler.Fields)
	public.POST("/knowledge", handler.CreateField)
}
