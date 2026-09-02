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

func (h *KnowledgeHandler) Domains(c *gin.Context) { result, err := h.service.SearchDomains(strings.TrimSpace(c.Query("q")), limit(c)); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error()); return }; utils.JSONResponse(c, http.StatusOK, "Knowledge domains retrieved successfully", result, nil) }
func (h *KnowledgeHandler) Subdomains(c *gin.Context) { id, err := uuid.Parse(c.Query("domain_id")); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "domain_id must be a valid UUID"); return }; result, err := h.service.SearchSubdomains(id, strings.TrimSpace(c.Query("q")), limit(c)); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error()); return }; utils.JSONResponse(c, http.StatusOK, "Knowledge subdomains retrieved successfully", result, nil) }
func (h *KnowledgeHandler) Fields(c *gin.Context) { id, err := uuid.Parse(c.Query("subdomain_id")); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusBadRequest, "subdomain_id must be a valid UUID"); return }; result, err := h.service.SearchFields(id, strings.TrimSpace(c.Query("q")), limit(c)); if err != nil { utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error()); return }; utils.JSONResponse(c, http.StatusOK, "Knowledge fields retrieved successfully", result, nil) }
