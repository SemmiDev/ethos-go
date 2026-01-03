package habit

import "time"

func (h *Habit) Update(name string, description *string, frequency Frequency, recurrence Recurrence, targetCount int, reminderTime *string) error {
	if name == "" {
		return ErrEmptyName
	}
	if targetCount < 1 {
		return ErrInvalidTargetCount
	}
	if err := frequency.Validate(); err != nil {
		return err
	}
	if err := recurrence.Validate(); err != nil {
		return err
	}
	if reminderTime != nil {
		if _, err := time.Parse("15:04", *reminderTime); err != nil {
			return ErrInvalidReminder
		}
	}

	h.name = name
	h.description = description
	h.frequency = frequency
	h.recurrence = recurrence
	h.targetCount = targetCount
	h.reminderTime = reminderTime
	h.updatedAt = time.Now()

	return nil
}
