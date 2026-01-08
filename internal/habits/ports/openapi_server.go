package ports

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/httputil"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/common/random"
	habits "github.com/semmidev/ethos-go/internal/generated/api/habits"
	"github.com/semmidev/ethos-go/internal/habits/app"
	"github.com/semmidev/ethos-go/internal/habits/app/command"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
)

type OpenAPIServer struct {
	app app.Application
}

func NewOpenAPIServer(app app.Application) *OpenAPIServer {
	return &OpenAPIServer{app: app}
}

// Ensure OpenAPIServer implements habits.ServerInterface
var _ habits.ServerInterface = (*OpenAPIServer)(nil)

// List all habits
// (GET /habits)
func (s *OpenAPIServer) ListHabits(w http.ResponseWriter, r *http.Request, params habits.ListHabitsParams) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	// Parse filter from query parameters
	filter := model.FilterFromRequest(r)

	// Handle legacy params for backward compatibility
	if params.Active != nil && *params.Active {
		active := true
		filter.IsActive = &active
	}
	if params.Inactive != nil && *params.Inactive {
		inactive := true
		filter.IsInactive = &inactive
	}

	result, err := s.app.Queries.ListHabits.Handle(r.Context(), query.ListHabits{
		UserID: user.UserID,
		Filter: filter,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	// Convert domain model to API model
	habitsList := make([]habits.Habit, 0, len(result.Habits))
	for _, h := range result.Habits {
		id, _ := uuid.Parse(h.HabitID)
		freq := habits.HabitFrequency(h.Frequency)
		habitsList = append(habitsList, habits.Habit{
			Id:           id,
			Name:         h.Name,
			Description:  h.Description,
			Frequency:    &freq,
			TargetCount:  &h.TargetCount,
			ReminderTime: h.ReminderTime,
			IsActive:     &h.IsActive,
			CreatedAt:    h.CreatedAt,
			UpdatedAt:    &h.UpdatedAt,
		})
	}

	// Return list with pagination in meta
	httputil.SuccessPaginated(w, r, habitsList, result.Pagination, "Habits retrieved successfully")
}

// Create a new habit
// (POST /habits)
func (s *OpenAPIServer) CreateHabit(w http.ResponseWriter, r *http.Request) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req habits.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	habitID := random.NewUUID().String()

	frequency := "daily"
	if req.Frequency != nil {
		frequency = *req.Frequency
	}

	targetCount := 1
	if req.TargetCount != nil {
		targetCount = *req.TargetCount
	}

	err = s.app.Commands.CreateHabit.Handle(r.Context(), command.CreateHabit{
		HabitID:      habitID,
		UserID:       user.UserID,
		Name:         req.Name,
		Description:  req.Description,
		Frequency:    frequency,
		TargetCount:  targetCount,
		ReminderTime: req.ReminderTime,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	// Fetch created habit
	h, err := s.app.Queries.GetHabit.Handle(r.Context(), query.GetHabit{
		HabitID: habitID,
		UserID:  user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	id, _ := uuid.Parse(h.HabitID)
	freq := habits.HabitFrequency(h.Frequency)
	resp := habits.Habit{
		Id:           id,
		Name:         h.Name,
		Description:  h.Description,
		Frequency:    &freq,
		TargetCount:  &h.TargetCount,
		ReminderTime: h.ReminderTime,
		IsActive:     &h.IsActive,
		CreatedAt:    h.CreatedAt,
		UpdatedAt:    &h.UpdatedAt,
	}

	httputil.Created(w, r, resp, "Habit created successfully")
}

// Delete a habit
// (DELETE /habits/{habitId})
func (s *OpenAPIServer) DeleteHabit(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.DeleteHabit.Handle(r.Context(), command.DeleteHabit{
		HabitID: habitId.String(),
		UserID:  user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Habit deleted successfully")
}

// Get a habit by ID
// (GET /habits/{habitId})
func (s *OpenAPIServer) GetHabit(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	h, err := s.app.Queries.GetHabit.Handle(r.Context(), query.GetHabit{
		HabitID: habitId.String(),
		UserID:  user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	id, _ := uuid.Parse(h.HabitID)
	freq := habits.HabitFrequency(h.Frequency)
	resp := habits.Habit{
		Id:           id,
		Name:         h.Name,
		Description:  h.Description,
		Frequency:    &freq,
		TargetCount:  &h.TargetCount,
		ReminderTime: h.ReminderTime,
		IsActive:     &h.IsActive,
		CreatedAt:    h.CreatedAt,
		UpdatedAt:    &h.UpdatedAt,
	}

	httputil.Success(w, r, resp, "Habit retrieved successfully")
}

// Update a habit
// (PUT /habits/{habitId})
func (s *OpenAPIServer) UpdateHabit(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req habits.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	err = s.app.Commands.UpdateHabit.Handle(r.Context(), command.UpdateHabit{
		HabitID:      habitId.String(),
		UserID:       user.UserID,
		Name:         req.Name,
		Description:  req.Description,
		Frequency:    req.Frequency,
		TargetCount:  req.TargetCount,
		ReminderTime: req.ReminderTime,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	// Fetch updated habit to return
	h, err := s.app.Queries.GetHabit.Handle(r.Context(), query.GetHabit{
		HabitID: habitId.String(),
		UserID:  user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	id, _ := uuid.Parse(h.HabitID)
	freq := habits.HabitFrequency(h.Frequency)
	resp := habits.Habit{
		Id:           id,
		Name:         h.Name,
		Description:  h.Description,
		Frequency:    &freq,
		TargetCount:  &h.TargetCount,
		ReminderTime: h.ReminderTime,
		IsActive:     &h.IsActive,
		CreatedAt:    h.CreatedAt,
		UpdatedAt:    &h.UpdatedAt,
	}

	httputil.Success(w, r, resp, "Habit updated successfully")
}

// Activate a habit
// (POST /habits/{habitId}/activate)
func (s *OpenAPIServer) ActivateHabit(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.ActivateHabit.Handle(r.Context(), command.ActivateHabit{
		HabitID: habitId.String(),
		UserID:  user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Habit activated successfully")
}

// Deactivate a habit
// (POST /habits/{habitId}/deactivate)
func (s *OpenAPIServer) DeactivateHabit(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.DeactivateHabit.Handle(r.Context(), command.DeactivateHabit{
		HabitID: habitId.String(),
		UserID:  user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Habit deactivated successfully")
}

// Get habit statistics
// (GET /habits/{habitId}/stats)
func (s *OpenAPIServer) GetHabitStats(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	stats, err := s.app.Queries.GetHabitStats.Handle(r.Context(), query.GetHabitStats{
		HabitID: habitId.String(),
		UserID:  user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	resp := habits.HabitStats{
		TotalLogs:     &stats.TotalCompletions,
		CurrentStreak: &stats.CurrentStreak,
		LongestStreak: &stats.LongestStreak,
	}

	httputil.Success(w, r, resp, "Habit stats retrieved successfully")
}

// Log a habit
// (POST /habits/{habitId}/logs)
func (s *OpenAPIServer) LogHabit(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req habits.LogHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	formattedDate := req.LogDate.Time

	logID := random.NewUUID().String()

	err = s.app.Commands.LogHabit.Handle(r.Context(), command.LogHabit{
		LogID:   logID,
		HabitID: habitId.String(),
		UserID:  user.UserID,
		LogDate: formattedDate,
		Count:   req.Count,
		Note:    req.Note, // Already *string from OpenAPI
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	// The spec returns object with log_id
	lID, _ := uuid.Parse(logID)
	resp := struct {
		LogId openapi_types.UUID `json:"log_id"`
	}{
		LogId: lID,
	}

	httputil.Created(w, r, resp, "Habit logged successfully")
}

// Get habit logs
// (GET /habits/{habitId}/logs)
func (s *OpenAPIServer) GetHabitLogs(w http.ResponseWriter, r *http.Request, habitId openapi_types.UUID, params habits.GetHabitLogsParams) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	// Parse filter from query parameters
	filter := model.FilterFromRequest(r)

	// Handle legacy date params for backward compatibility
	if params.StartDate != nil {
		t := params.StartDate.Time
		filter.StartDate = &t
	}
	if params.EndDate != nil {
		t := params.EndDate.Time
		filter.EndDate = &t
	}

	result, err := s.app.Queries.GetHabitLogs.Handle(r.Context(), query.GetHabitLogs{
		HabitID: habitId.String(),
		UserID:  user.UserID,
		Filter:  filter,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	logsList := make([]habits.HabitLog, 0, len(result.Logs))
	for _, l := range result.Logs {
		lID, _ := uuid.Parse(l.LogID)
		hID, _ := uuid.Parse(l.HabitID)

		d := openapi_types.Date{Time: l.LogDate}

		logsList = append(logsList, habits.HabitLog{
			Id:        &lID,
			HabitId:   &hID,
			LogDate:   &d,
			Count:     &l.Count,
			Note:      l.Note, // Already *string
			CreatedAt: &l.CreatedAt,
		})
	}

	// Return list with pagination in meta
	httputil.SuccessPaginated(w, r, logsList, result.Pagination, "Habit logs retrieved successfully")
}

// Update a habit log
// (PUT /habit-logs/{logId})
func (s *OpenAPIServer) UpdateHabitLog(w http.ResponseWriter, r *http.Request, logId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req habits.UpdateHabitLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	var logDate *time.Time
	if req.LogDate != nil {
		logDate = req.LogDate
	}

	err = s.app.Commands.UpdateHabitLog.Handle(r.Context(), command.UpdateHabitLog{
		LogID:   logId.String(),
		UserID:  user.UserID,
		Count:   req.Count,
		Note:    req.Note,
		LogDate: logDate,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Habit log updated successfully")
}

// Delete a habit log
// (DELETE /habit-logs/{logId})
func (s *OpenAPIServer) DeleteHabitLog(w http.ResponseWriter, r *http.Request, logId openapi_types.UUID) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.DeleteHabitLog.Handle(r.Context(), command.DeleteHabitLog{
		LogID:  logId.String(),
		UserID: user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Habit log deleted successfully")
}

// Get user dashboard
// (GET /dashboard)
func (s *OpenAPIServer) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	dashboard, err := s.app.Queries.GetDashboard.Handle(r.Context(), query.GetDashboard{
		UserID: user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	resp := habits.Dashboard{
		ActiveHabitsCount: &dashboard.TotalActiveHabits,
		TotalLogsToday:    &dashboard.TotalCompletionsToday,
		CurrentStreak:     &dashboard.CurrentStreak,
		LongestStreak:     &dashboard.LongestStreak,
		WeeklyCompletion:  &dashboard.WeeklyCompletion,
		TotalLogs:         &dashboard.TotalLogs,
	}

	httputil.Success(w, r, resp, "Dashboard data retrieved successfully")
}

// Get weekly analytics data
// (GET /analytics/weekly)
func (s *OpenAPIServer) GetWeeklyAnalytics(w http.ResponseWriter, r *http.Request) {
	user, err := authctx.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	analytics, err := s.app.Queries.GetWeeklyAnalytics.Handle(r.Context(), query.GetWeeklyAnalytics{
		UserID: user.UserID,
	})

	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	// Convert to OpenAPI response type
	days := make([]habits.DailyAnalytics, len(analytics.Days))
	for i, day := range analytics.Days {
		dayName := day.DayName
		date := openapi_types.Date{}
		_ = date.UnmarshalText([]byte(day.Date))
		logsCount := day.LogsCount
		completion := day.CompletionPercentage
		days[i] = habits.DailyAnalytics{
			DayName:              &dayName,
			Date:                 &date,
			LogsCount:            &logsCount,
			CompletionPercentage: &completion,
		}
	}

	avgCompletion := analytics.AverageCompletion
	resp := habits.WeeklyAnalytics{
		Days:              &days,
		AverageCompletion: &avgCompletion,
	}

	httputil.Success(w, r, resp, "Weekly analytics retrieved successfully")
}
