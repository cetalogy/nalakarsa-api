package utils

import (
	"strconv"

	"nalakarsa/internal/dto"

	"github.com/gin-gonic/gin"
)

// JSONResponse sends a successful JSON response
func JSONResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination *dto.PaginationResponse) {
	c.JSON(statusCode, dto.APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// ErrorJSONResponse sends an error JSON response
func ErrorJSONResponse(c *gin.Context, statusCode int, message string, errors []string) {
	c.JSON(statusCode, dto.APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// ErrorJSONResponseWithMessage sends a simple single-message error JSON response
func ErrorJSONResponseWithMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, dto.APIResponse{
		Success: false,
		Message: message,
	})
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

	return page, limit
}
