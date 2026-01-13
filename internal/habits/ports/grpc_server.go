package ports

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
	"github.com/semmidev/ethos-go/internal/common/grpcutil"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/common/random"
	commonv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/common/v1"
	habitsv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/habits/v1"
	"github.com/semmidev/ethos-go/internal/habits/app"
	"github.com/semmidev/ethos-go/internal/habits/app/command"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
)

// HabitsGRPCServer implements the gRPC HabitsService interface.
type HabitsGRPCServer struct {
	habitsv1.UnimplementedHabitsServiceServer
	app app.Application
}

// NewHabitsGRPCServer creates a new HabitsGRPCServer.
func NewHabitsGRPCServer(application app.Application) *HabitsGRPCServer {
	return &HabitsGRPCServer{app: application}
}

// ListHabits returns all habits for the authenticated user.
func (s *HabitsGRPCServer) ListHabits(ctx context.Context, req *habitsv1.ListHabitsRequest) (*habitsv1.ListHabitsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	filter := model.NewFilter()
	if req.Page > 0 {
		filter.CurrentPage = int(req.Page)
	}
	if req.PerPage > 0 {
		filter.PerPage = int(req.PerPage)
	}
	if req.Active != nil && *req.Active {
		active := true
		filter.IsActive = &active
	}
	if req.Inactive != nil && *req.Inactive {
		inactive := true
		filter.IsInactive = &inactive
	}
	if req.Keyword != nil {
		filter.Keyword = *req.Keyword
	}
	if req.SortBy != nil {
		filter.SortBy = *req.SortBy
	}
	if req.SortDirection != nil {
		filter.SortDirection = *req.SortDirection
	}

	result, err := s.app.Queries.ListHabits.Handle(ctx, query.ListHabits{
		UserID: user.UserID,
		Filter: filter,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	habits := make([]*habitsv1.Habit, 0, len(result.Habits))
	for _, h := range result.Habits {
		habits = append(habits, toProtoHabit(h))
	}

	return &habitsv1.ListHabitsResponse{
		Success: true,
		Message: "Habits retrieved successfully",
		Data:    habits,
		Meta: &commonv1.Meta{
			Pagination: &commonv1.PaginationResponse{
				HasPreviousPage:        result.Pagination.HasPreviousPage,
				HasNextPage:            result.Pagination.HasNextPage,
				CurrentPage:            int32(result.Pagination.CurrentPage),
				PerPage:                int32(result.Pagination.PerPage),
				TotalData:              int32(result.Pagination.TotalData),
				TotalDataInCurrentPage: int32(result.Pagination.TotalDataInCurrentPage),
				LastPage:               int32(result.Pagination.LastPage),
				From:                   int32(result.Pagination.From),
				To:                     int32(result.Pagination.To),
			},
		},
	}, nil
}

// CreateHabit creates a new habit.
func (s *HabitsGRPCServer) CreateHabit(ctx context.Context, req *habitsv1.CreateHabitRequest) (*habitsv1.HabitResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	habitID := random.NewUUID().String()

	frequency := "daily"
	if req.Frequency != nil {
		frequency = *req.Frequency
	}

	targetCount := 1
	if req.TargetCount != nil {
		targetCount = int(*req.TargetCount)
	}

	cmd := command.CreateHabit{
		HabitID:      habitID,
		UserID:       user.UserID,
		Name:         req.Name,
		Description:  req.Description,
		Frequency:    frequency,
		TargetCount:  targetCount,
		ReminderTime: req.ReminderTime,
	}

	if err := s.app.Commands.CreateHabit.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	h, err := s.app.Queries.GetHabit.Handle(ctx, query.GetHabit{
		HabitID: habitID,
		UserID:  user.UserID,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.HabitResponse{
		Success: true,
		Message: "Habit created successfully",
		Data:    toProtoHabit(*h),
	}, nil
}

// GetHabit retrieves a habit by ID.
func (s *HabitsGRPCServer) GetHabit(ctx context.Context, req *habitsv1.GetHabitRequest) (*habitsv1.HabitResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	h, err := s.app.Queries.GetHabit.Handle(ctx, query.GetHabit{
		HabitID: req.HabitId,
		UserID:  user.UserID,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.HabitResponse{
		Success: true,
		Message: "Habit retrieved successfully",
		Data:    toProtoHabit(*h),
	}, nil
}

// UpdateHabit updates a habit.
func (s *HabitsGRPCServer) UpdateHabit(ctx context.Context, req *habitsv1.UpdateHabitRequest) (*habitsv1.HabitResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	var targetCount *int
	if req.TargetCount != nil {
		tc := int(*req.TargetCount)
		targetCount = &tc
	}

	cmd := command.UpdateHabit{
		HabitID:      req.HabitId,
		UserID:       user.UserID,
		Name:         req.Name,
		Description:  req.Description,
		Frequency:    req.Frequency,
		TargetCount:  targetCount,
		ReminderTime: req.ReminderTime,
	}

	if err := s.app.Commands.UpdateHabit.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	h, err := s.app.Queries.GetHabit.Handle(ctx, query.GetHabit{
		HabitID: req.HabitId,
		UserID:  user.UserID,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.HabitResponse{
		Success: true,
		Message: "Habit updated successfully",
		Data:    toProtoHabit(*h),
	}, nil
}

// DeleteHabit deletes a habit.
func (s *HabitsGRPCServer) DeleteHabit(ctx context.Context, req *habitsv1.DeleteHabitRequest) (*habitsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.DeleteHabit{
		HabitID: req.HabitId,
		UserID:  user.UserID,
	}

	if err := s.app.Commands.DeleteHabit.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.SuccessResponse{
		Success: true,
		Message: "Habit deleted successfully",
	}, nil
}

// ActivateHabit activates a habit.
func (s *HabitsGRPCServer) ActivateHabit(ctx context.Context, req *habitsv1.ActivateHabitRequest) (*habitsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.ActivateHabit{
		HabitID: req.HabitId,
		UserID:  user.UserID,
	}

	if err := s.app.Commands.ActivateHabit.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.SuccessResponse{
		Success: true,
		Message: "Habit activated successfully",
	}, nil
}

// DeactivateHabit deactivates a habit.
func (s *HabitsGRPCServer) DeactivateHabit(ctx context.Context, req *habitsv1.DeactivateHabitRequest) (*habitsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.DeactivateHabit{
		HabitID: req.HabitId,
		UserID:  user.UserID,
	}

	if err := s.app.Commands.DeactivateHabit.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.SuccessResponse{
		Success: true,
		Message: "Habit deactivated successfully",
	}, nil
}

// GetHabitStats retrieves habit statistics.
func (s *HabitsGRPCServer) GetHabitStats(ctx context.Context, req *habitsv1.GetHabitStatsRequest) (*habitsv1.HabitStatsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	stats, err := s.app.Queries.GetHabitStats.Handle(ctx, query.GetHabitStats{
		HabitID: req.HabitId,
		UserID:  user.UserID,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.HabitStatsResponse{
		Success: true,
		Message: "Habit stats retrieved successfully",
		Data: &habitsv1.HabitStats{
			TotalLogs:     int32(stats.TotalCompletions),
			CurrentStreak: int32(stats.CurrentStreak),
			LongestStreak: int32(stats.LongestStreak),
		},
	}, nil
}

// LogHabit logs a habit completion.
func (s *HabitsGRPCServer) LogHabit(ctx context.Context, req *habitsv1.LogHabitRequest) (*habitsv1.LogHabitResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	logDate, err := time.Parse("2006-01-02", req.LogDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid log_date format, expected YYYY-MM-DD")
	}

	logID := random.NewUUID().String()

	cmd := command.LogHabit{
		LogID:   logID,
		HabitID: req.HabitId,
		UserID:  user.UserID,
		LogDate: logDate,
		Count:   int(req.Count),
		Note:    req.Note,
	}

	if err := s.app.Commands.LogHabit.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.LogHabitResponse{
		Success: true,
		Message: "Habit logged successfully",
		Data: &habitsv1.LogHabitData{
			LogId: logID,
		},
	}, nil
}

// GetHabitLogs retrieves logs for a habit.
func (s *HabitsGRPCServer) GetHabitLogs(ctx context.Context, req *habitsv1.GetHabitLogsRequest) (*habitsv1.GetHabitLogsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	filter := model.NewFilter()
	if req.Page > 0 {
		filter.CurrentPage = int(req.Page)
	}
	if req.PerPage > 0 {
		filter.PerPage = int(req.PerPage)
	}
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		filter.StartDate = &t
	}
	if req.EndDate != nil {
		t, _ := time.Parse("2006-01-02", *req.EndDate)
		filter.EndDate = &t
	}
	if req.Keyword != nil {
		filter.Keyword = *req.Keyword
	}

	result, err := s.app.Queries.GetHabitLogs.Handle(ctx, query.GetHabitLogs{
		HabitID: req.HabitId,
		UserID:  user.UserID,
		Filter:  filter,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	logs := make([]*habitsv1.HabitLog, 0, len(result.Logs))
	for _, l := range result.Logs {
		logs = append(logs, &habitsv1.HabitLog{
			Id:        l.LogID,
			HabitId:   l.HabitID,
			LogDate:   l.LogDate.Format("2006-01-02"),
			Count:     int32(l.Count),
			Note:      l.Note,
			CreatedAt: timestamppb.New(l.CreatedAt),
		})
	}

	return &habitsv1.GetHabitLogsResponse{
		Success: true,
		Message: "Habit logs retrieved successfully",
		Data:    logs,
		Meta: &commonv1.Meta{
			Pagination: &commonv1.PaginationResponse{
				HasPreviousPage:        result.Pagination.HasPreviousPage,
				HasNextPage:            result.Pagination.HasNextPage,
				CurrentPage:            int32(result.Pagination.CurrentPage),
				PerPage:                int32(result.Pagination.PerPage),
				TotalData:              int32(result.Pagination.TotalData),
				TotalDataInCurrentPage: int32(result.Pagination.TotalDataInCurrentPage),
				LastPage:               int32(result.Pagination.LastPage),
				From:                   int32(result.Pagination.From),
				To:                     int32(result.Pagination.To),
			},
		},
	}, nil
}

// UpdateHabitLog updates a habit log.
func (s *HabitsGRPCServer) UpdateHabitLog(ctx context.Context, req *habitsv1.UpdateHabitLogRequest) (*habitsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	var count *int
	if req.Count != nil {
		c := int(*req.Count)
		count = &c
	}

	var logDate *time.Time
	if req.LogDate != nil {
		t, err := time.Parse("2006-01-02", *req.LogDate)
		if err == nil {
			logDate = &t
		}
	}

	cmd := command.UpdateHabitLog{
		LogID:   req.LogId,
		UserID:  user.UserID,
		Count:   count,
		Note:    req.Note,
		LogDate: logDate,
	}

	if err := s.app.Commands.UpdateHabitLog.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.SuccessResponse{
		Success: true,
		Message: "Habit log updated successfully",
	}, nil
}

// DeleteHabitLog deletes a habit log.
func (s *HabitsGRPCServer) DeleteHabitLog(ctx context.Context, req *habitsv1.DeleteHabitLogRequest) (*habitsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.DeleteHabitLog{
		LogID:  req.LogId,
		UserID: user.UserID,
	}

	if err := s.app.Commands.DeleteHabitLog.Handle(ctx, cmd); err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.SuccessResponse{
		Success: true,
		Message: "Habit log deleted successfully",
	}, nil
}

// GetDashboard retrieves the user's dashboard data.
func (s *HabitsGRPCServer) GetDashboard(ctx context.Context, req *habitsv1.GetDashboardRequest) (*habitsv1.DashboardResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	dashboard, err := s.app.Queries.GetDashboard.Handle(ctx, query.GetDashboard{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	return &habitsv1.DashboardResponse{
		Success: true,
		Message: "Dashboard data retrieved successfully",
		Data: &habitsv1.Dashboard{
			ActiveHabitsCount: int32(dashboard.TotalActiveHabits),
			TotalLogsToday:    int32(dashboard.TotalCompletionsToday),
			CurrentStreak:     int32(dashboard.CurrentStreak),
			LongestStreak:     int32(dashboard.LongestStreak),
			WeeklyCompletion:  int32(dashboard.WeeklyCompletion),
			TotalLogs:         int32(dashboard.TotalLogs),
		},
	}, nil
}

// GetWeeklyAnalytics retrieves weekly analytics data.
func (s *HabitsGRPCServer) GetWeeklyAnalytics(ctx context.Context, req *habitsv1.GetWeeklyAnalyticsRequest) (*habitsv1.WeeklyAnalyticsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	analytics, err := s.app.Queries.GetWeeklyAnalytics.Handle(ctx, query.GetWeeklyAnalytics{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, toHabitsGRPCError(err)
	}

	days := make([]*habitsv1.DailyAnalytics, len(analytics.Days))
	for i, day := range analytics.Days {
		days[i] = &habitsv1.DailyAnalytics{
			DayName:              day.DayName,
			Date:                 day.Date,
			LogsCount:            int32(day.LogsCount),
			CompletionPercentage: int32(day.CompletionPercentage),
		}
	}

	return &habitsv1.WeeklyAnalyticsResponse{
		Success: true,
		Message: "Weekly analytics retrieved successfully",
		Data: &habitsv1.WeeklyAnalytics{
			Days:              days,
			AverageCompletion: int32(analytics.AverageCompletion),
		},
	}, nil
}

// toProtoHabit converts a query.Habit to a protobuf Habit.
func toProtoHabit(h query.Habit) *habitsv1.Habit {
	habit := &habitsv1.Habit{
		Id:          h.HabitID,
		Name:        h.Name,
		Frequency:   h.Frequency,
		TargetCount: int32(h.TargetCount),
		IsActive:    h.IsActive,
		CreatedAt:   timestamppb.New(h.CreatedAt),
		UpdatedAt:   timestamppb.New(h.UpdatedAt),
	}

	if h.Description != nil {
		habit.Description = h.Description
	}
	if h.ReminderTime != nil {
		habit.ReminderTime = h.ReminderTime
	}

	return habit
}

// toHabitsGRPCError converts application errors to gRPC status errors.
func toHabitsGRPCError(err error) error {
	return grpcutil.ToGRPCError(err)
}
