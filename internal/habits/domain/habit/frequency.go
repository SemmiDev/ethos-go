package habit

import "errors"

type Frequency struct {
	value string
}

const (
	FrequencyDaily   = "daily"
	FrequencyWeekly  = "weekly"
	FrequencyMonthly = "monthly"
)

func NewFrequency(value string) (Frequency, error) {
	f := Frequency{value: value}
	if err := f.Validate(); err != nil {
		return Frequency{}, err
	}
	return f, nil
}

func (f Frequency) Validate() error {
	switch f.value {
	case FrequencyDaily, FrequencyWeekly, FrequencyMonthly:
		return nil
	default:
		return errors.New("invalid frequency: must be daily, weekly, or monthly")
	}
}

func (f Frequency) String() string {
	return f.value
}

func (f Frequency) IsDaily() bool   { return f.value == FrequencyDaily }
func (f Frequency) IsWeekly() bool  { return f.value == FrequencyWeekly }
func (f Frequency) IsMonthly() bool { return f.value == FrequencyMonthly }
