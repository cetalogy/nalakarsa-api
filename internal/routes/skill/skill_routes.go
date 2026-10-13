package skill

import (
	skillhandler "nalakarsa/internal/handler/skill"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public, protected *gin.RouterGroup, handler *skillhandler.SkillHandler) {
	public.GET("/skills/search", handler.Search)
	public.POST("/skills", handler.Create)
	public.GET("/skills/users/:user_id", handler.ListUserSkills)
	protected.GET("/users/me/skills", handler.ListMySkills)
	protected.POST("/users/me/skills", handler.AddMySkill)
	protected.DELETE("/users/me/skills/:skill_id", handler.RemoveMySkill)
}
