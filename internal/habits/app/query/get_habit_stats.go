package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// GetHabitStats query retrieves statistics for a habit
type GetHabitStats struct {
	HabitID string
	UserID  string
}

// GetHabitStatsHandler processes get habit stats queries
type GetHabitStatsHandler decorator.QueryHandler[GetHabitStats, *HabitStats]

// GetHabitStatsReadModel interface for data access
type GetHabitStatsReadModel interface {
	GetHabitStats(ctx context.Context, habitID, userID string) (*HabitStats, error)
}

type getHabitStatsHandler struct {
	readModel GetHabitStatsReadModel
}

// NewGetHabitStatsHandler creates a new handler with decorators
func NewGetHabitStatsHandler(
	readModel GetHabitStatsReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetHabitStatsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators(
		getHabitStatsHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h getHabitStatsHandler) Handle(ctx context.Context, q GetHabitStats) (*HabitStats, error) {
	return h.readModel.GetHabitStats(ctx, q.HabitID, q.UserID)
}
