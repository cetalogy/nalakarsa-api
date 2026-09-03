package knowledgeroutes

import (
	knowledgehandler "nalakarsa/internal/handler/knowledge"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, handler *knowledgehandler.KnowledgeHandler) {
	public.GET("/knowledge/fields/search", handler.Fields)
}
