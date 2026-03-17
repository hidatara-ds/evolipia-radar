package api

import (
	"github.com/gin-gonic/gin"
)

// APIErrorCode defines standard internal error codes for the API gateway.
type APIErrorCode string

const (
	ErrCodeValidation   APIErrorCode = "VALIDATION_ERROR"
	ErrCodeInternal     APIErrorCode = "INTERNAL_ERROR"
	ErrCodeTimeout      APIErrorCode = "TIMEOUT_ERROR"
	ErrCodeUnauthorized APIErrorCode = "UNAUTHORIZED"
)

// APIError represents the standard JSON error payload structure.
type APIError struct {
	Message string       `json:"message"`
	Code    APIErrorCode `json:"code"`
}

// errorResponse wraps the error payload to match the {"error": ...} JSON structure.
type errorResponse struct {
	Error APIError `json:"error"`
}

// RespondWithError is a helper to abort the request and return standardized JSON.
func RespondWithError(c *gin.Context, status int, code APIErrorCode, message string) {
	c.AbortWithStatusJSON(status, errorResponse{
		Error: APIError{
			Message: message,
			Code:    code,
		},
	})
}
