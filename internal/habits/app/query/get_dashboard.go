package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// GetDashboard query retrieves dashboard summary for a user
type GetDashboard struct {
	UserID string
}

// GetDashboardHandler processes get dashboard queries
type GetDashboardHandler decorator.QueryHandler[GetDashboard, *DashboardSummary]

// GetDashboardReadModel interface for data access
type GetDashboardReadModel interface {
	GetDashboard(ctx context.Context, userID string) (*DashboardSummary, error)
}

type getDashboardHandler struct {
	readModel GetDashboardReadModel
}

// NewGetDashboardHandler creates a new handler with decorators
func NewGetDashboardHandler(
	readModel GetDashboardReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetDashboardHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators(
		getDashboardHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h getDashboardHandler) Handle(ctx context.Context, q GetDashboard) (*DashboardSummary, error) {
	return h.readModel.GetDashboard(ctx, q.UserID)
}
