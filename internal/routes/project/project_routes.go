package projectroutes

import (
	"nalakarsa/internal/handler/project"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *projecthandler.ProjectHandler) {
	projects := public.Group("/projects")
	projects.GET("", h.List)
	projects.GET("/:id", h.GetByID)
	projects.GET("/:id/members", h.ListMembers)

	protected.POST("/projects", h.Create)
	protected.PATCH("/projects/:id", h.Update)
	protected.DELETE("/projects/:id", h.Delete)
	protected.POST("/projects/:id/applications", h.Apply)
	protected.GET("/projects/:id/applications", h.ListApplications)
	protected.PATCH("/projects/:id/applications/:applicationId", h.UpdateApplicationStatus)
	protected.POST("/projects/:id/milestones", h.CreateMilestone)
	protected.PATCH("/projects/:id/milestones/:milestoneId", h.UpdateMilestone)
}
