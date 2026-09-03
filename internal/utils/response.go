package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"nalakarsa/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
func JSONResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination *dto.PaginationResponse) {
	c.JSON(statusCode, dto.APIResponse{
		Data:      data,
		Meta:      pagination,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
func ErrorJSONResponse(c *gin.Context, statusCode int, code string, message string, details map[string]string) {
	c.JSON(statusCode, dto.APIErrorResponse{
		Error: dto.APIErrorDetail{
			Code:      code,
			Message:   message,
			Details:   details,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}
func ErrorJSONResponseWithMessage(c *gin.Context, statusCode int, message string) {
	code := "INTERNAL_ERROR"
	details := make(map[string]string)

	switch statusCode {
	case 400:
		code = "BAD_REQUEST"
		details["reason"] = message
		details["suggestion"] = "Verify query parameters, path variables, or request format"
	case 401:
		code = "UNAUTHORIZED"
		details["reason"] = message
		msgLower := strings.ToLower(message)
		if strings.Contains(msgLower, "email") || strings.Contains(msgLower, "password") || strings.Contains(msgLower, "credential") {
			details["suggestion"] = "Please check your email address and password, or register a new account"
		} else {
			details["suggestion"] = "Provide a valid JWT token in Authorization header or perform a login/token refresh"
		}
	case 403:
		code = "FORBIDDEN"
		details["reason"] = message
		details["suggestion"] = "Ensure you are logged in with the correct account or resource ownership"
	case 404:
		code = "NOT_FOUND"
		details["reason"] = message
		details["suggestion"] = "Check the requested ID or endpoint path to ensure resource exists"
	case 409:
		code = "CONFLICT"
		details["reason"] = message
		details["suggestion"] = "The request conflicts with existing data or resource state"
	case 429:
		code = "RATE_LIMITED"
		details["reason"] = message
		details["suggestion"] = "Too many requests. Please wait a moment before trying again"
	default:
		code = "INTERNAL_ERROR"
		details["reason"] = message
		details["suggestion"] = "An unexpected server condition occurred. Please try again or contact support"
	}

	ErrorJSONResponse(c, statusCode, code, message, details)
}
func ValidationErrorResponse(c *gin.Context, input interface{}) {
	var details map[string]string
	switch v := input.(type) {
	case map[string]string:
		details = v
	case error:
		details = FormatBindingError(v)
	}
	ErrorJSONResponse(c, 400, "VALIDATION_ERROR", "Validation failed", details)
}
func FormatBindingError(err error) map[string]string {
	details := make(map[string]string)
	if err == nil {
		return details
	}
	if errors.Is(err, io.EOF) || err.Error() == "EOF" {
		details["body"] = "Request body cannot be empty. Please provide a valid JSON payload."
		return details
	}
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		details["body"] = fmt.Sprintf("Malformed JSON syntax at offset %d", syntaxErr.Offset)
		return details
	}
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		fieldName := strings.ToLower(unmarshalErr.Field)
		if fieldName == "" {
			fieldName = "body"
		}
		details[fieldName] = fmt.Sprintf("Expected %s type, but received %s", unmarshalErr.Type.String(), unmarshalErr.Value)
		return details
	}
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				details[field] = fmt.Sprintf("%s is required", field)
			case "email":
				details[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				details[field] = fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
			case "max":
				details[field] = fmt.Sprintf("%s cannot exceed %s characters", field, fe.Param())
			case "oneof":
				details[field] = fmt.Sprintf("%s must be one of [%s]", field, fe.Param())
			default:
				details[field] = fmt.Sprintf("%s failed on validation rule '%s'", field, fe.Tag())
			}
		}
		return details
	}
	details["body"] = fmt.Sprintf("Invalid request payload: %s", err.Error())
	return details
}
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
