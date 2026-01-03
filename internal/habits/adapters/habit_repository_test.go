//go:build integration

package adapters_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/habits/adapters"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
	"github.com/semmidev/ethos-go/internal/testutil"
)

var testContainer *testutil.PostgresContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = testutil.NewPostgresContainer(ctx)
	if err != nil {
		panic("failed to start test container: " + err.Error())
	}

	if err := testContainer.RunMigrations(ctx); err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	exitCode := m.Run()

	testContainer.Cleanup(ctx)
	if exitCode != 0 {
		panic("tests failed")
	}
}

func TestHabitRepository_AddAndGet(t *testing.T) {
	ctx := context.Background()
	testContainer.TruncateTables(ctx)

	// Create test user
	user := testutil.DefaultTestUser()
	if err := testContainer.CreateTestUser(ctx, user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := adapters.NewHabitPostgresRepository(testContainer.DB)

	// Create habit
	freq, _ := habit.NewFrequency("daily")
	h, err := habit.NewHabit(
		uuid.New().String(),
		user.UserID.String(),
		"Test Habit",
		nil,
		freq,
		habit.DefaultRecurrence(),
		1,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	// Add habit
	if err := repo.AddHabit(ctx, h); err != nil {
		t.Fatalf("failed to add habit: %v", err)
	}

	// Get habit back
	retrieved, err := repo.GetHabit(ctx, h.HabitID(), user.UserID.String())
	if err != nil {
		t.Fatalf("failed to get habit: %v", err)
	}

	if retrieved.HabitID() != h.HabitID() {
		t.Errorf("expected habit ID %s, got %s", h.HabitID(), retrieved.HabitID())
	}
	if retrieved.Name() != "Test Habit" {
		t.Errorf("expected name 'Test Habit', got %s", retrieved.Name())
	}
}

func TestHabitRepository_GetHabit_NotFound(t *testing.T) {
	ctx := context.Background()
	testContainer.TruncateTables(ctx)

	user := testutil.DefaultTestUser()
	if err := testContainer.CreateTestUser(ctx, user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := adapters.NewHabitPostgresRepository(testContainer.DB)

	_, err := repo.GetHabit(ctx, uuid.New().String(), user.UserID.String())
	if err == nil {
		t.Error("expected error for non-existent habit")
	}
}

func TestHabitRepository_GetHabit_Unauthorized(t *testing.T) {
	ctx := context.Background()
	testContainer.TruncateTables(ctx)

	owner := testutil.DefaultTestUser()
	owner.Email = "owner@example.com"
	if err := testContainer.CreateTestUser(ctx, owner); err != nil {
		t.Fatalf("failed to create owner: %v", err)
	}

	other := testutil.TestUser{
		UserID: uuid.New(),
		Email:  "other@example.com",
		Name:   "Other User",
	}
	if err := testContainer.CreateTestUser(ctx, other); err != nil {
		t.Fatalf("failed to create other user: %v", err)
	}

	repo := adapters.NewHabitPostgresRepository(testContainer.DB)

	// Create habit for owner
	freq, _ := habit.NewFrequency("daily")
	h, _ := habit.NewHabit(uuid.New().String(), owner.UserID.String(), "Owner Habit", nil, freq, habit.DefaultRecurrence(), 1, nil)
	repo.AddHabit(ctx, h)

	// Try to access as other user
	_, err := repo.GetHabit(ctx, h.HabitID(), other.UserID.String())
	if err == nil {
		t.Error("expected unauthorized error when accessing another user's habit")
	}
}

func TestHabitRepository_UpdateHabit(t *testing.T) {
	ctx := context.Background()
	testContainer.TruncateTables(ctx)

	user := testutil.DefaultTestUser()
	if err := testContainer.CreateTestUser(ctx, user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := adapters.NewHabitPostgresRepository(testContainer.DB)

	freq, _ := habit.NewFrequency("daily")
	h, _ := habit.NewHabit(uuid.New().String(), user.UserID.String(), "Original Name", nil, freq, habit.DefaultRecurrence(), 1, nil)
	repo.AddHabit(ctx, h)

	// Update habit
	err := repo.UpdateHabit(ctx, h.HabitID(), user.UserID.String(), func(ctx context.Context, existing *habit.Habit) (*habit.Habit, error) {
		newDesc := "Updated description"
		err := existing.Update("Updated Name", &newDesc, existing.Frequency(), existing.Recurrence(), existing.TargetCount(), existing.ReminderTime())
		return existing, err
	})
	if err != nil {
		t.Fatalf("failed to update habit: %v", err)
	}

	// Verify update
	updated, _ := repo.GetHabit(ctx, h.HabitID(), user.UserID.String())
	if updated.Name() != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %s", updated.Name())
	}
}

func TestHabitRepository_DeleteHabit(t *testing.T) {
	ctx := context.Background()
	testContainer.TruncateTables(ctx)

	user := testutil.DefaultTestUser()
	if err := testContainer.CreateTestUser(ctx, user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := adapters.NewHabitPostgresRepository(testContainer.DB)

	freq, _ := habit.NewFrequency("daily")
	h, _ := habit.NewHabit(uuid.New().String(), user.UserID.String(), "To Delete", nil, freq, habit.DefaultRecurrence(), 1, nil)
	repo.AddHabit(ctx, h)

	// Delete
	err := repo.DeleteHabit(ctx, h.HabitID(), user.UserID.String())
	if err != nil {
		t.Fatalf("failed to delete habit: %v", err)
	}

	// Verify deletion
	_, err = repo.GetHabit(ctx, h.HabitID(), user.UserID.String())
	if err == nil {
		t.Error("expected not found error after deletion")
	}
}

func TestHabitRepository_ListHabitsByUser(t *testing.T) {
	ctx := context.Background()
	testContainer.TruncateTables(ctx)

	user := testutil.DefaultTestUser()
	if err := testContainer.CreateTestUser(ctx, user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := adapters.NewHabitPostgresRepository(testContainer.DB)

	// Create multiple habits
	freq, _ := habit.NewFrequency("daily")
	for i := 0; i < 3; i++ {
		h, _ := habit.NewHabit(uuid.New().String(), user.UserID.String(), "Habit "+string(rune('A'+i)), nil, freq, habit.DefaultRecurrence(), 1, nil)
		repo.AddHabit(ctx, h)
	}

	habits, err := repo.ListHabitsByUser(ctx, user.UserID.String())
	if err != nil {
		t.Fatalf("failed to list habits: %v", err)
	}

	if len(habits) != 3 {
		t.Errorf("expected 3 habits, got %d", len(habits))
	}
}
