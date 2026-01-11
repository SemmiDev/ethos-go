package grpcutil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// CustomHTTPError is a custom error handler for gRPC-Gateway
func CustomHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st, _ := status.FromError(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	resp := StandardResponse{
		Success: false,
		Message: st.Message(),
	}

	// Extract error details
	errorResponse := map[string]interface{}{
		"code":    st.Code().String(), // Use the gRPC code name for now, e.g. "INVALID_ARGUMENT"
		"message": st.Message(),
	}

	// Check for details in the status
	if details := st.Details(); len(details) > 0 {
		// We expect the first detail to be our structpb.Struct map
		if s, ok := details[0].(*structpb.Struct); ok {
			errorResponse["details"] = s.AsMap()
		}
	}

	resp.Error = errorResponse

	// Use encoding/json to marshal the response
	enc := json.NewEncoder(w)
	_ = enc.Encode(resp)
}
