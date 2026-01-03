package habit_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

func TestFrequency(t *testing.T) {
	t.Parallel()

	Convey("Given the Frequency value object", t, func() {

		Convey("When creating with valid values", func() {
			validFrequencies := []string{"daily", "weekly", "monthly"}

			for _, freq := range validFrequencies {
				freq := freq
				Convey("Then '"+freq+"' should be valid", func() {
					f, err := habit.NewFrequency(freq)
					So(err, ShouldBeNil)
					So(f.String(), ShouldEqual, freq)
				})
			}
		})

		Convey("When creating with invalid values", func() {
			invalidValues := []string{"", "hourly", "yearly", "custom", "DAILY", "Daily"}

			for _, val := range invalidValues {
				val := val
				name := val
				if name == "" {
					name = "empty"
				}
				Convey("Then '"+name+"' should be invalid", func() {
					_, err := habit.NewFrequency(val)
					So(err, ShouldNotBeNil)
				})
			}
		})
	})
}

func TestFrequencyMethods(t *testing.T) {
	t.Parallel()

	Convey("Given frequency type checks", t, func() {

		Convey("When checking IsDaily", func() {
			daily, _ := habit.NewFrequency("daily")
			weekly, _ := habit.NewFrequency("weekly")

			Convey("Then daily frequency should return true", func() {
				So(daily.IsDaily(), ShouldBeTrue)
			})

			Convey("Then weekly frequency should return false", func() {
				So(weekly.IsDaily(), ShouldBeFalse)
			})
		})

		Convey("When checking IsWeekly", func() {
			weekly, _ := habit.NewFrequency("weekly")
			daily, _ := habit.NewFrequency("daily")

			Convey("Then weekly frequency should return true", func() {
				So(weekly.IsWeekly(), ShouldBeTrue)
			})

			Convey("Then daily frequency should return false", func() {
				So(daily.IsWeekly(), ShouldBeFalse)
			})
		})

		Convey("When checking IsMonthly", func() {
			monthly, _ := habit.NewFrequency("monthly")
			daily, _ := habit.NewFrequency("daily")

			Convey("Then monthly frequency should return true", func() {
				So(monthly.IsMonthly(), ShouldBeTrue)
			})

			Convey("Then daily frequency should return false", func() {
				So(daily.IsMonthly(), ShouldBeFalse)
			})
		})

		Convey("When calling String()", func() {
			freq, _ := habit.NewFrequency("daily")

			Convey("Then it should return the frequency value", func() {
				So(freq.String(), ShouldEqual, "daily")
			})
		})
	})
}
