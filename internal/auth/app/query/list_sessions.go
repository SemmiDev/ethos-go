package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/model"
)

// ListSessionsQuery requests all sessions for a user.
type ListSessionsQuery struct {
	UserID           string
	CurrentSessionID string       // To mark which session is the current one
	IncludeBlocked   bool         // Whether to include blocked sessions
	IncludeExpired   bool         // Whether to include expired sessions
	Filter           model.Filter // Pagination and sorting
}

// ListSessionsResult contains the paginated list of sessions
type ListSessionsResult struct {
	Sessions   []*SessionDTO
	Pagination *model.Paging
}

// ListSessionsHandler retrieves all sessions for a user.
type ListSessionsHandler decorator.QueryHandler[ListSessionsQuery, ListSessionsResult]

// ListSessionsReadModel interface for data access
type ListSessionsReadModel interface {
	ListSessions(ctx context.Context, userID uuid.UUID, includeBlocked, includeExpired bool, filter model.Filter) ([]*session.Session, int, error)
}

type listSessionsHandler struct {
	readModel ListSessionsReadModel
}

// NewListSessionsHandler creates a handler with its dependencies.
func NewListSessionsHandler(
	readModel ListSessionsReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ListSessionsHandler {
	return decorator.ApplyQueryDecorators[ListSessionsQuery, ListSessionsResult](
		listSessionsHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

// Handle executes the query to list all user sessions.
func (h listSessionsHandler) Handle(ctx context.Context, query ListSessionsQuery) (ListSessionsResult, error) {
	// Parse the user ID
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return ListSessionsResult{}, apperror.InvalidInput("user_id", "invalid UUID format").
			WithError(err)
	}

	// Validate filter
	query.Filter.Validate()

	// Fetch sessions via read model
	sessions, totalCount, err := h.readModel.ListSessions(ctx, userID, query.IncludeBlocked, query.IncludeExpired, query.Filter)
	if err != nil {
		return ListSessionsResult{}, h.translateError(err, "list sessions")
	}

	// Convert domain models to DTOs
	dtos := make([]*SessionDTO, len(sessions))
	for i, s := range sessions {
		isCurrent := s.SessionID.String() == query.CurrentSessionID
		dtos[i] = toSessionDTO(s, isCurrent)
	}

	// Create pagination info
	pagination, err := model.NewPaging(query.Filter.CurrentPage, query.Filter.PerPage, totalCount)
	if err != nil {
		return ListSessionsResult{}, err
	}

	return ListSessionsResult{
		Sessions:   dtos,
		Pagination: pagination,
	}, nil
}

// translateError converts domain errors to AppErrors
func (h listSessionsHandler) translateError(err error, operation string) *apperror.AppError {
	// Domain errors translation if any

	if appErr := apperror.GetAppError(err); appErr != nil {
		return appErr
	}

	return apperror.DatabaseError(operation, err)
}

// toSessionDTO converts a domain Session to a SessionDTO.
// This helper function encapsulates the conversion logic so we don't
// repeat it in multiple places.
func toSessionDTO(s *session.Session, isCurrent bool) *SessionDTO {
	return &SessionDTO{
		SessionID: s.SessionID.String(),
		UserAgent: s.UserAgent,
		ClientIP:  s.ClientIP,
		IsBlocked: s.IsBlocked,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
		IsActive:  s.IsValid(), // Use domain logic to compute this
		IsCurrent: isCurrent,
	}
}
