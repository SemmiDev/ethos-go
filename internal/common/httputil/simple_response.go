package httputil

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/model"
)

// StandardResponse is the unified response structure
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ResponseMeta contains metadata for responses (e.g., pagination)
type ResponseMeta struct {
	Pagination *model.Paging `json:"pagination,omitempty"`
}

// Success processes a successful request and returns a JSON response
func Success(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// SuccessWithMeta processes a successful request with metadata (e.g., pagination)
func SuccessWithMeta(w http.ResponseWriter, r *http.Request, data interface{}, meta *ResponseMeta, message string) {
	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// SuccessPaginated is a helper for returning paginated list responses
// Data should be the list items, pagination will be placed in meta field
func SuccessPaginated(w http.ResponseWriter, r *http.Request, data interface{}, pagination *model.Paging, message string) {
	meta := &ResponseMeta{
		Pagination: pagination,
	}
	SuccessWithMeta(w, r, data, meta, message)
}

// Created processes a creation request and returns a 201 JSON response
func Created(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

// Error processes an error and returns a JSON error response
func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Default values
	statusCode := http.StatusInternalServerError
	resp := StandardResponse{
		Success: false,
		Message: "Internal Server Error",
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		statusCode = appErr.HTTPStatusCode()
		// If message is set in AppError, use it.
		// NOTE: appErr.Error() usually combines Message + underlying Err.
		// For the response "message" field, we might prefer just appErr.Message which is user-safe.
		// But appErr.Error() might include internal error info if configured so?
		// apperror.go: Error() returns fmt.Sprintf("%s: %v", e.Message, e.Err) if Err != nil
		// apperror.go: Message is "human-readable error message safe to show to clients"

		// Let's use Message for the top level message field
		resp.Message = appErr.Message

		errData := map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
		}

		if len(appErr.Details) > 0 {
			errData["details"] = appErr.Details
		}

		resp.Error = errData
	} else {
		// Handle domain errors by inspecting error message patterns
		errMsg := err.Error()
		errMsgLower := strings.ToLower(errMsg)

		switch {
		case strings.Contains(errMsgLower, "not found"):
			statusCode = http.StatusNotFound
			resp.Message = capitalizeFirst(errMsg)
			resp.Error = map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": capitalizeFirst(errMsg),
			}
		case strings.Contains(errMsgLower, "unauthorized") ||
			strings.Contains(errMsgLower, "cannot access"):
			statusCode = http.StatusForbidden
			resp.Message = capitalizeFirst(errMsg)
			resp.Error = map[string]interface{}{
				"code":    "FORBIDDEN",
				"message": capitalizeFirst(errMsg),
			}
		case strings.Contains(errMsgLower, "already"):
			statusCode = http.StatusConflict
			resp.Message = capitalizeFirst(errMsg)
			resp.Error = map[string]interface{}{
				"code":    "CONFLICT",
				"message": capitalizeFirst(errMsg),
			}
		case strings.Contains(errMsgLower, "invalid") ||
			strings.Contains(errMsgLower, "empty") ||
			strings.Contains(errMsgLower, "must be"):
			statusCode = http.StatusBadRequest
			resp.Message = capitalizeFirst(errMsg)
			resp.Error = map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": capitalizeFirst(errMsg),
			}
		default:
			// Generic internal error - don't expose internal error details
			resp.Error = map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			}
		}
	}

	render.Status(r, statusCode)
	render.JSON(w, r, resp)
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
