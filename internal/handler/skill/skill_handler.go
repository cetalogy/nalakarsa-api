package skill

import (
	"net/http"
	"strconv"

	"nalakarsa/internal/dto"
	skillservice "nalakarsa/internal/service/skill"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SkillHandler struct{ service skillservice.Service }

func NewSkillHandler(service skillservice.Service) *SkillHandler {
	return &SkillHandler{service: service}
}

func currentUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := value.(uuid.UUID)
	return id, ok
}

func (h *SkillHandler) Search(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := h.service.Search(c.Query("q"), limit)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Skills retrieved successfully", items, nil)
}

func (h *SkillHandler) Create(c *gin.Context) {
	var req dto.CreateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	item, err := h.service.Create(req.Name)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Skill created successfully", item, nil)
}

func (h *SkillHandler) ListMySkills(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	h.listUserSkills(c, userID)
}

func (h *SkillHandler) ListUserSkills(c *gin.Context) {
	id, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	h.listUserSkills(c, id)
}

func (h *SkillHandler) listUserSkills(c *gin.Context, userID uuid.UUID) {
	items, err := h.service.ListByUser(userID)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "User skills retrieved successfully", items, nil)
}

func (h *SkillHandler) AddMySkill(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req dto.AddUserSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	item, err := h.service.AddToUser(userID, req.SkillID, req.Name)
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "Skill added to profile successfully", item, nil)
}

func (h *SkillHandler) RemoveMySkill(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	skillID, err := uuid.Parse(c.Param("skill_id"))
	if err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.RemoveFromUser(userID, skillID); err != nil {
		utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, "Skill removed from profile successfully", nil, nil)
}
