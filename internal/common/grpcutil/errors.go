package grpcutil

import (
	"net/http"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// ToGRPCError converts application errors to gRPC status errors with rich details.
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	appErr := apperror.GetAppError(err)
	if appErr == nil {
		// Generic internal error
		return status.Error(codes.Internal, err.Error())
	}

	code := toGRPCCode(appErr.StatusCode)
	st := status.New(code, appErr.Message)

	details := appErr.Details
	if appErr.Code != "" {
		if details == nil {
			details = make(map[string]interface{})
		}
		// Create a copy to avoid mutating the original error's details if referenced elsewhere
		newDetails := make(map[string]interface{}, len(details)+1)
		for k, v := range details {
			newDetails[k] = v
		}
		newDetails["_code"] = appErr.Code
		details = newDetails
	}

	if len(details) > 0 {
		// Convert details map to structpb.Struct
		detailsStruct, err := structpb.NewStruct(details)
		if err == nil {
			// Attach the details struct to the status
			st, _ = st.WithDetails(detailsStruct)
		}
	}

	return st.Err()
}

func toGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		return codes.FailedPrecondition
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusInternalServerError:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
