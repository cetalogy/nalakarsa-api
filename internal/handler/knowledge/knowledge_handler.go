package knowledgehandler

import (
	"net/http"
	"strconv"
	"strings"

	knowledgeservice "nalakarsa/internal/service/knowledge"
	"nalakarsa/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KnowledgeHandler struct{ service knowledgeservice.KnowledgeService }
func NewKnowledgeHandler(service knowledgeservice.KnowledgeService) *KnowledgeHandler { return &KnowledgeHandler{service: service} }
func limit(c *gin.Context) int { n, err := strconv.Atoi(c.Query("limit")); if err != nil || n <= 0 { return 10 }; if n > 30 { return 30 }; return n }

func (h *KnowledgeHandler) Fields(c *gin.Context) { id := uuid.Nil; if value := strings.TrimSpace(c.Query("subdomain_id")); value != "" { parsed, err := uuid.Parse(value); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "subdomain_id must be a valid UUID"); return }; id = parsed }; result, err := h.service.SearchFields(id, strings.TrimSpace(c.Query("q")), limit(c)); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error()); return }; utils.JSONResponse(c, http.StatusOK, "Knowledge fields retrieved successfully", result, nil) }
