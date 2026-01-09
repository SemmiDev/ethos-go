package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// GetHabit query retrieves a single habit by ID
type GetHabit struct {
	HabitID string
	UserID  string
}

// GetHabitHandler processes get habit queries
type GetHabitHandler decorator.QueryHandler[GetHabit, *Habit]

// GetHabitReadModel interface for data access
type GetHabitReadModel interface {
	GetHabitQuery(ctx context.Context, habitID, userID string) (*Habit, error)
}

type getHabitHandler struct {
	readModel GetHabitReadModel
}

// NewGetHabitHandler creates a new handler with decorators
func NewGetHabitHandler(
	readModel GetHabitReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetHabitHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators(
		getHabitHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h getHabitHandler) Handle(ctx context.Context, q GetHabit) (*Habit, error) {
	return h.readModel.GetHabitQuery(ctx, q.HabitID, q.UserID)
}
