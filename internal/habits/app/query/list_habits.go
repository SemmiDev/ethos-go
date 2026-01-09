package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/model"
)

// ListHabits query retrieves all habits for a user with filtering and pagination
type ListHabits struct {
	UserID string
	Filter model.Filter
}

// ListHabitsResult contains the paginated list of habits
type ListHabitsResult struct {
	Habits     []Habit
	Pagination *model.Paging
}

// ListHabitsHandler processes list habits queries
type ListHabitsHandler decorator.QueryHandler[ListHabits, ListHabitsResult]

// ListHabitsReadModel interface for data access
type ListHabitsReadModel interface {
	ListHabits(ctx context.Context, userID string, filter model.Filter) ([]Habit, int, error)
}

type listHabitsHandler struct {
	readModel ListHabitsReadModel
}

// NewListHabitsHandler creates a new handler with decorators
func NewListHabitsHandler(
	readModel ListHabitsReadModel,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ListHabitsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return decorator.ApplyQueryDecorators(
		listHabitsHandler{readModel: readModel},
		log,
		metricsClient,
	)
}

func (h listHabitsHandler) Handle(ctx context.Context, q ListHabits) (ListHabitsResult, error) {
	// Validate filter
	q.Filter.Validate()

	// Validate allowed sort columns
	allowedSortColumns := []string{"name", "created_at", "updated_at", "is_active"}
	q.Filter.ValidateSortBy(allowedSortColumns)

	habits, totalCount, err := h.readModel.ListHabits(ctx, q.UserID, q.Filter)
	if err != nil {
		return ListHabitsResult{}, err
	}

	pagination, err := model.NewPaging(q.Filter.CurrentPage, q.Filter.PerPage, totalCount)
	if err != nil {
		return ListHabitsResult{}, err
	}

	return ListHabitsResult{
		Habits:     habits,
		Pagination: pagination,
	}, nil
}
