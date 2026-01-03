package habit

import (
	"errors"
	"time"
)

// Recurrence represents advanced scheduling configuration for a habit
type Recurrence struct {
	days     int16 // Bitmask: Sun=1, Mon=2, Tue=4, Wed=8, Thu=16, Fri=32, Sat=64
	interval int   // Every N periods (days/weeks/months based on frequency)
}

// Day constants for bitmask
const (
	Sunday    int16 = 1 << iota // 1
	Monday                      // 2
	Tuesday                     // 4
	Wednesday                   // 8
	Thursday                    // 16
	Friday                      // 32
	Saturday                    // 64
	AllDays   int16 = 127       // All days (1+2+4+8+16+32+64)
	Weekdays  int16 = 62        // Mon-Fri (2+4+8+16+32)
	Weekends  int16 = 65        // Sat-Sun (1+64)
)

// NewRecurrence creates a new Recurrence with validation
func NewRecurrence(days int16, interval int) (Recurrence, error) {
	r := Recurrence{days: days, interval: interval}
	if err := r.Validate(); err != nil {
		return Recurrence{}, err
	}
	return r, nil
}

// Validate checks if the recurrence is valid
func (r Recurrence) Validate() error {
	if r.days < 0 || r.days > AllDays {
		return errors.New("invalid recurrence days: must be between 0 and 127")
	}
	if r.interval < 1 {
		return errors.New("invalid recurrence interval: must be at least 1")
	}
	return nil
}

// DefaultRecurrence returns a recurrence for every day
func DefaultRecurrence() Recurrence {
	return Recurrence{days: AllDays, interval: 1}
}

// Getters
func (r Recurrence) Days() int16   { return r.days }
func (r Recurrence) Interval() int { return r.interval }

// HasDay checks if a specific day is included in the recurrence
func (r Recurrence) HasDay(day int16) bool {
	return r.days&day != 0
}

// ShouldCompleteOn determines if a habit should be completed on a given date
// based on the frequency, recurrence days, and interval
func (r Recurrence) ShouldCompleteOn(date time.Time, frequency Frequency, habitCreatedAt time.Time) bool {
	weekday := date.Weekday()
	dayBit := int16(1 << weekday)

	// Check if this day of the week is included
	if r.days != 0 && !r.HasDay(dayBit) {
		return false
	}

	// For daily frequency with interval > 1, check if today is an interval day
	if frequency.IsDaily() && r.interval > 1 {
		daysSinceCreation := int(date.Sub(habitCreatedAt).Hours() / 24)
		if daysSinceCreation%r.interval != 0 {
			return false
		}
	}

	// For weekly frequency with interval > 1
	if frequency.IsWeekly() && r.interval > 1 {
		_, createdWeek := habitCreatedAt.ISOWeek()
		_, currentWeek := date.ISOWeek()
		weeksDiff := currentWeek - createdWeek
		if weeksDiff%r.interval != 0 {
			return false
		}
	}

	// For monthly frequency with interval > 1
	if frequency.IsMonthly() && r.interval > 1 {
		monthsDiff := (date.Year()-habitCreatedAt.Year())*12 + int(date.Month()) - int(habitCreatedAt.Month())
		if monthsDiff%r.interval != 0 {
			return false
		}
	}

	return true
}

// DayNames returns a slice of day names included in this recurrence
func (r Recurrence) DayNames() []string {
	days := []string{}
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, name := range dayNames {
		if r.HasDay(1 << i) {
			days = append(days, name)
		}
	}
	return days
}

// IsEveryDay returns true if all days are selected
func (r Recurrence) IsEveryDay() bool {
	return r.days == AllDays
}
