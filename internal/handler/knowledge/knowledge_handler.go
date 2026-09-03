package knowledgehandler

import (
	"net/http"
	"strconv"
	"strings"

	"nalakarsa/internal/dto"
	knowledgeservice "nalakarsa/internal/service/knowledge"
	"nalakarsa/internal/utils"
	"github.com/gin-gonic/gin"
)

type KnowledgeHandler struct{ service knowledgeservice.KnowledgeService }
func NewKnowledgeHandler(service knowledgeservice.KnowledgeService) *KnowledgeHandler { return &KnowledgeHandler{service: service} }
func limit(c *gin.Context) int { n, err := strconv.Atoi(c.Query("limit")); if err != nil || n <= 0 { return 10 }; if n > 30 { return 30 }; return n }

func (h *KnowledgeHandler) Fields(c *gin.Context) { result, err := h.service.SearchFields(strings.TrimSpace(c.Query("q")), limit(c)); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error()); return }; utils.JSONResponse(c, http.StatusOK, "Knowledge fields retrieved successfully", result, nil) }
func (h *KnowledgeHandler) CreateField(c *gin.Context) { var req dto.CreateReferenceRequest; if err := c.ShouldBindJSON(&req); err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "name is required"); return }; item, err := h.service.CreateField(req.Name); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, err.Error()); return }; utils.JSONResponse(c, http.StatusCreated, "Knowledge field created successfully", item, nil) }
