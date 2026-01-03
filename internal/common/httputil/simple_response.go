package httputil

import (
	"errors"
	"net/http"

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
		// Generic error
		resp.Error = map[string]interface{}{
			"message": resp.Message, // Internal Server Error
		}
	}

	render.Status(r, statusCode)
	render.JSON(w, r, resp)
}
