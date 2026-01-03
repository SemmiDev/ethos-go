package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// GetWeeklyAnalytics query retrieves weekly analytics for a user
type GetWeeklyAnalytics struct {
	UserID string
}

// GetWeeklyAnalyticsHandler processes weekly analytics queries
type GetWeeklyAnalyticsHandler decorator.QueryHandler[GetWeeklyAnalytics, *WeeklyAnalytics]

// GetWeeklyAnalyticsReadModel interface for data access
type GetWeeklyAnalyticsReadModel interface {
	GetWeeklyAnalytics(ctx context.Context, userID string) (*WeeklyAnalytics, error)
}

type getWeeklyAnalyticsHandler struct {
	readModel GetWeeklyAnalyticsReadModel
}

// NewGetWeeklyAnalyticsHandler creates a new handler with decorators
func NewGetWeeklyAnalyticsHandler(
	readModel GetWeeklyAnalyticsReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetWeeklyAnalyticsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators[GetWeeklyAnalytics, *WeeklyAnalytics](
		getWeeklyAnalyticsHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h getWeeklyAnalyticsHandler) Handle(ctx context.Context, q GetWeeklyAnalytics) (*WeeklyAnalytics, error) {
	return h.readModel.GetWeeklyAnalytics(ctx, q.UserID)
}
