package knowledgeroutes

import (
	"github.com/gin-gonic/gin"
	knowledgehandler "nalakarsa/internal/handler/knowledge"
)

func RegisterRoutes(public *gin.RouterGroup, handler *knowledgehandler.KnowledgeHandler) {
	public.GET("/knowledge/search", handler.Fields)
	public.POST("/knowledge", handler.CreateField)
}
