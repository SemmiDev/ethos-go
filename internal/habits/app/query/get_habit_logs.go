package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/model"
)

// GetHabitLogs query retrieves logs for a habit with filtering and pagination
type GetHabitLogs struct {
	HabitID string
	UserID  string
	Filter  model.Filter
}

// GetHabitLogsResult contains the paginated list of habit logs
type GetHabitLogsResult struct {
	Logs       []HabitLog
	Pagination *model.Paging
}

// GetHabitLogsHandler processes get habit logs queries
type GetHabitLogsHandler decorator.QueryHandler[GetHabitLogs, GetHabitLogsResult]

// GetHabitLogsReadModel interface for data access
type GetHabitLogsReadModel interface {
	GetHabitLogs(ctx context.Context, habitID, userID string, filter model.Filter) ([]HabitLog, int, error)
}

type getHabitLogsHandler struct {
	readModel GetHabitLogsReadModel
}

// NewGetHabitLogsHandler creates a new handler with decorators
func NewGetHabitLogsHandler(
	readModel GetHabitLogsReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetHabitLogsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators(
		getHabitLogsHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h getHabitLogsHandler) Handle(ctx context.Context, q GetHabitLogs) (GetHabitLogsResult, error) {
	// Validate filter
	q.Filter.Validate()

	// Validate allowed sort columns
	allowedSortColumns := []string{"log_date", "created_at", "count"}
	q.Filter.ValidateSortBy(allowedSortColumns)

	logs, totalCount, err := h.readModel.GetHabitLogs(ctx, q.HabitID, q.UserID, q.Filter)
	if err != nil {
		return GetHabitLogsResult{}, err
	}

	pagination, err := model.NewPaging(q.Filter.CurrentPage, q.Filter.PerPage, totalCount)
	if err != nil {
		return GetHabitLogsResult{}, err
	}

	return GetHabitLogsResult{
		Logs:       logs,
		Pagination: pagination,
	}, nil
}
