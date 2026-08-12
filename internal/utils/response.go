package utils

import (
	"strconv"

	"nalakarsa/internal/dto"

	"github.com/gin-gonic/gin"
)

// JSONResponse sends a successful JSON response using the standard envelope:
// { "data": ..., "meta": ..., "message": "..." }
func JSONResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination *dto.PaginationResponse) {
	c.JSON(statusCode, dto.APIResponse{
		Data:    data,
		Meta:    pagination,
		Message: message,
	})
}

// ErrorJSONResponse sends an error JSON response with error code and optional details:
// { "error": { "code": "...", "message": "...", "details": { ... } } }
func ErrorJSONResponse(c *gin.Context, statusCode int, code string, message string, details map[string]string) {
	c.JSON(statusCode, dto.APIErrorResponse{
		Error: dto.APIErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ErrorJSONResponseWithMessage sends a simple error response with just code and message
func ErrorJSONResponseWithMessage(c *gin.Context, statusCode int, message string) {
	code := "INTERNAL_ERROR"
	switch {
	case statusCode == 400:
		code = "BAD_REQUEST"
	case statusCode == 401:
		code = "UNAUTHORIZED"
	case statusCode == 403:
		code = "FORBIDDEN"
	case statusCode == 404:
		code = "NOT_FOUND"
	case statusCode == 409:
		code = "CONFLICT"
	case statusCode == 429:
		code = "RATE_LIMITED"
	}
	ErrorJSONResponse(c, statusCode, code, message, nil)
}

// ValidationErrorResponse sends a validation error response with per-field details
func ValidationErrorResponse(c *gin.Context, details map[string]string) {
	ErrorJSONResponse(c, 400, "VALIDATION_ERROR", "Data tidak valid", details)
}

// ParsePaginationRequest helper to parse page and limit queries
func ParsePaginationRequest(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit
}
