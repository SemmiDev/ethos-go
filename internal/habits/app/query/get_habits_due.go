package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

type GetHabitsDue struct{}

type GetHabitsDueHandler decorator.QueryHandler[GetHabitsDue, []ReminderHabit]

type HabitsDueReadModel interface {
	GetHabitsDueForReminder(ctx context.Context) ([]ReminderHabit, error)
}

type getHabitsDueHandler struct {
	readModel HabitsDueReadModel
}

func NewGetHabitsDueHandler(
	readModel HabitsDueReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetHabitsDueHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators[GetHabitsDue, []ReminderHabit](
		getHabitsDueHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h getHabitsDueHandler) Handle(ctx context.Context, _ GetHabitsDue) ([]ReminderHabit, error) {
	return h.readModel.GetHabitsDueForReminder(ctx)
}
