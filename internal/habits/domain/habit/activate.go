package habit

import "time"

func (h *Habit) Activate() error {
	if h.isActive {
		return ErrAlreadyActive
	}

	h.isActive = true
	h.updatedAt = time.Now()
	return nil
}
