package habit_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

func TestRecurrence(t *testing.T) {
	Convey("Given the Recurrence value object", t, func() {

		Convey("When creating default recurrence", func() {
			r := habit.DefaultRecurrence()

			Convey("Then days should be 127 (all days)", func() {
				So(r.Days(), ShouldEqual, 127)
			})

			Convey("Then interval should be 1", func() {
				So(r.Interval(), ShouldEqual, 1)
			})
		})

		Convey("When creating with valid parameters", func() {
			// Monday + Wednesday + Friday = 2 + 8 + 32 = 42
			r, err := habit.NewRecurrence(42, 2)

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then days should be 42", func() {
				So(r.Days(), ShouldEqual, 42)
			})

			Convey("Then interval should be 2", func() {
				So(r.Interval(), ShouldEqual, 2)
			})
		})

		Convey("When creating with invalid interval", func() {
			testCases := []struct {
				name     string
				interval int
			}{
				{"zero interval", 0},
				{"negative interval", -1},
			}

			for _, tc := range testCases {
				tc := tc
				Convey("Then "+tc.name+" should return an error", func() {
					_, err := habit.NewRecurrence(127, tc.interval)
					So(err, ShouldNotBeNil)
				})
			}
		})

		Convey("When creating with valid interval", func() {
			testCases := []struct {
				name     string
				interval int
			}{
				{"interval 1", 1},
				{"large interval 30", 30},
			}

			for _, tc := range testCases {
				tc := tc
				Convey("Then "+tc.name+" should succeed", func() {
					_, err := habit.NewRecurrence(127, tc.interval)
					So(err, ShouldBeNil)
				})
			}
		})
	})
}

func TestRecurrenceHasDay(t *testing.T) {
	Convey("Given recurrence day checking", t, func() {

		Convey("When checking Monday-only recurrence", func() {
			// Monday = bit 1 = 2
			mondayOnly, _ := habit.NewRecurrence(2, 1)

			// Bitmask values: Sun=1, Mon=2, Tue=4, Wed=8, Thu=16, Fri=32, Sat=64
			Convey("Then Monday (bit 2) should be included", func() {
				So(mondayOnly.HasDay(2), ShouldBeTrue)
			})

			Convey("Then Tuesday (bit 4) should NOT be included", func() {
				So(mondayOnly.HasDay(4), ShouldBeFalse)
			})

			Convey("Then Sunday (bit 1) should NOT be included", func() {
				So(mondayOnly.HasDay(1), ShouldBeFalse)
			})
		})

		Convey("When checking all-days recurrence", func() {
			allDays, _ := habit.NewRecurrence(127, 1)

			// All days bitmask: 1+2+4+8+16+32+64 = 127
			dayBits := []int16{1, 2, 4, 8, 16, 32, 64}
			dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

			for i, bit := range dayBits {
				bit := bit
				name := dayNames[i]
				Convey("Then "+name+" should be included", func() {
					So(allDays.HasDay(bit), ShouldBeTrue)
				})
			}
		})
	})
}

func TestRecurrenceShouldCompleteOn(t *testing.T) {
	Convey("Given habit completion scheduling", t, func() {

		Convey("When daily habit with interval 1", func() {
			r, _ := habit.NewRecurrence(127, 1) // all days, interval 1
			freq, _ := habit.NewFrequency("daily")
			createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			checkDate := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)

			Convey("Then it should complete every day", func() {
				So(r.ShouldCompleteOn(checkDate, freq, createdAt), ShouldBeTrue)
			})
		})

		Convey("When daily habit with interval 3", func() {
			r, _ := habit.NewRecurrence(127, 3) // all days, every 3 days
			freq, _ := habit.NewFrequency("daily")
			createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

			testCases := []struct {
				date   time.Time
				expect bool
			}{
				{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), true},  // Day 0 (creation day)
				{time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), false}, // Day 1
				{time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC), false}, // Day 2
				{time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC), true},  // Day 3
				{time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC), true},  // Day 6
			}

			for _, tc := range testCases {
				tc := tc
				label := tc.date.Format("Jan 02")
				if tc.expect {
					Convey("Then "+label+" should require completion", func() {
						So(r.ShouldCompleteOn(tc.date, freq, createdAt), ShouldBeTrue)
					})
				} else {
					Convey("Then "+label+" should NOT require completion", func() {
						So(r.ShouldCompleteOn(tc.date, freq, createdAt), ShouldBeFalse)
					})
				}
			}
		})

		Convey("When weekly habit with specific days (Mon + Fri)", func() {
			// Monday and Friday = 2 + 32 = 34
			r, _ := habit.NewRecurrence(34, 1)
			freq, _ := habit.NewFrequency("weekly")
			createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) // Wednesday

			testCases := []struct {
				date   time.Time
				day    string
				expect bool
			}{
				{time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC), "Monday", true},
				{time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC), "Tuesday", false},
				{time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC), "Friday", true},
				{time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), "Saturday", false},
			}

			for _, tc := range testCases {
				tc := tc
				if tc.expect {
					Convey(fmt.Sprintf("Then %s should require completion", tc.day), func() {
						So(r.ShouldCompleteOn(tc.date, freq, createdAt), ShouldBeTrue)
					})
				} else {
					Convey(fmt.Sprintf("Then %s should NOT require completion", tc.day), func() {
						So(r.ShouldCompleteOn(tc.date, freq, createdAt), ShouldBeFalse)
					})
				}
			}
		})
	})
}
