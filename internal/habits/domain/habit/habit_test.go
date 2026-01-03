package habit_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

func TestNewHabit(t *testing.T) {
	t.Parallel()

	Convey("Given the Habit domain", t, func() {

		Convey("When creating a habit with valid input", func() {
			desc := "Morning exercise routine"
			freq, _ := habit.NewFrequency("daily")
			recurrence := habit.DefaultRecurrence()

			h, err := habit.NewHabit(
				"habit-123",
				"user-456",
				"Exercise",
				&desc,
				freq,
				recurrence,
				3,
				nil,
			)

			Convey("Then it should succeed without error", func() {
				So(err, ShouldBeNil)
				So(h, ShouldNotBeNil)
			})

			Convey("Then it should have correct HabitID", func() {
				So(h.HabitID(), ShouldEqual, "habit-123")
			})

			Convey("Then it should have correct UserID", func() {
				So(h.UserID(), ShouldEqual, "user-456")
			})

			Convey("Then it should have correct Name", func() {
				So(h.Name(), ShouldEqual, "Exercise")
			})

			Convey("Then it should have correct Description", func() {
				So(*h.Description(), ShouldEqual, desc)
			})

			Convey("Then it should have correct TargetCount", func() {
				So(h.TargetCount(), ShouldEqual, 3)
			})

			Convey("Then it should be active by default", func() {
				So(h.IsActive(), ShouldBeTrue)
			})
		})
	})
}

func TestNewHabitValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		habitID     string
		userID      string
		habitName   string
		targetCount int
		shouldError bool
	}{
		{"empty habit ID", "", "user-123", "Test", 1, true},
		{"empty user ID", "habit-123", "", "Test", 1, true},
		{"empty name", "habit-123", "user-456", "", 1, true},
		{"zero target", "habit-123", "user-456", "Test", 0, true},
		{"negative target", "habit-123", "user-456", "Test", -1, true},
		{"valid input", "habit-123", "user-456", "Test", 1, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			Convey("Given "+tc.name, t, func() {
				freq, _ := habit.NewFrequency("daily")
				_, err := habit.NewHabit(
					tc.habitID, tc.userID, tc.habitName, nil,
					freq, habit.DefaultRecurrence(), tc.targetCount, nil,
				)

				if tc.shouldError {
					Convey("Then it should return an error", func() {
						So(err, ShouldNotBeNil)
					})
				} else {
					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
					})
				}
			})
		})
	}
}

func TestReminderTimeValidation(t *testing.T) {
	t.Parallel()

	Convey("Given reminder time validation", t, func() {
		freq, _ := habit.NewFrequency("daily")

		Convey("When reminder time is invalid format", func() {
			invalidTime := "25:99"
			_, err := habit.NewHabit("h-1", "u-1", "Test", nil, freq, habit.DefaultRecurrence(), 1, &invalidTime)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When reminder time is valid format", func() {
			validTime := "08:30"
			h, err := habit.NewHabit("h-1", "u-1", "Test", nil, freq, habit.DefaultRecurrence(), 1, &validTime)

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
				So(*h.ReminderTime(), ShouldEqual, validTime)
			})
		})
	})
}

func TestHabitAuthorization(t *testing.T) {
	t.Parallel()

	Convey("Given a habit owned by a user", t, func() {
		freq, _ := habit.NewFrequency("daily")
		h, _ := habit.NewHabit("h-1", "user-owner", "Test", nil, freq, habit.DefaultRecurrence(), 1, nil)

		Convey("When the owner tries to view", func() {
			err := h.CanBeViewedBy("user-owner")

			Convey("Then it should allow access", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When a non-owner tries to view", func() {
			err := h.CanBeViewedBy("user-other")

			Convey("Then it should deny access", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUnmarshalHabitFromDatabase(t *testing.T) {
	t.Parallel()

	Convey("Given database habit data", t, func() {
		desc := "test description"
		now := time.Now()

		h, err := habit.UnmarshalHabitFromDatabase(
			"habit-db-1",
			"user-db-1",
			"DB Habit",
			&desc,
			"weekly",
			127, // all days
			1,
			2,
			nil,
			true,
			now,
			now,
		)

		Convey("Then it should unmarshal successfully", func() {
			So(err, ShouldBeNil)
			So(h, ShouldNotBeNil)
		})

		Convey("Then it should have correct HabitID", func() {
			So(h.HabitID(), ShouldEqual, "habit-db-1")
		})

		Convey("Then it should have correct Frequency", func() {
			So(h.Frequency().String(), ShouldEqual, "weekly")
		})
	})
}
