package habit

import "time"

func (h *Habit) Deactivate() error {
	if !h.isActive {
		return ErrAlreadyInactive
	}

	h.isActive = false
	h.updatedAt = time.Now()
	return nil
}
