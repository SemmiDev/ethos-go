import { useState, useEffect } from 'react';
import { Modal } from '../ui/Modal';
import { Input, TextArea, Select } from '../ui/Input';
import { Button } from '../ui/Button';
import { useHabitsStore } from '../../stores/habitsStore';
import { useUIStore } from '../../stores/uiStore';

const frequencyOptions = [
    { value: 'daily', label: 'Daily' },
    { value: 'weekly', label: 'Weekly' },
    { value: 'monthly', label: 'Monthly' },
];

const DAYS_OF_WEEK = [
    { bit: 1, label: 'Su' },
    { bit: 2, label: 'Mo' },
    { bit: 4, label: 'Tu' },
    { bit: 8, label: 'We' },
    { bit: 16, label: 'Th' },
    { bit: 32, label: 'Fr' },
    { bit: 64, label: 'Sa' },
];

export function CreateHabitModal({ isOpen, onClose }) {
    const { createHabit, isLoading } = useHabitsStore();
    const { addToast } = useUIStore();
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        frequency: 'daily',
        target_count: 1,
        reminder_time: '',
        recurrence_days: 127, // Default all days
        recurrence_interval: 1,
    });
    const [errors, setErrors] = useState({});

    const validate = () => {
        const newErrors = {};
        if (!formData.name.trim()) newErrors.name = 'Name is required';
        if (formData.target_count < 1) newErrors.target_count = 'Target must be at least 1';
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!validate()) return;

        const submitData = getSubmitData();
        const result = await createHabit({
            ...submitData,
            target_count: parseInt(formData.target_count, 10),
            recurrence_interval: parseInt(formData.recurrence_interval, 10),
        });

        if (result.success) {
            addToast({ type: 'success', title: 'Habit Created', message: 'Your new habit has been created successfully.' });
            setFormData({ name: '', description: '', frequency: 'daily', target_count: 1, reminder_time: '' });
            onClose();
        } else {
            addToast({ type: 'error', title: 'Error', message: result.error });
        }
    };

    const handleChange = (field) => (e) => {
        setFormData({ ...formData, [field]: e.target.value });
        if (errors[field]) setErrors({ ...errors, [field]: null });
    };

    // Helper for passing null instead of empty string for optional time
    const getSubmitData = () => {
        const data = { ...formData };
        if (!data.reminder_time) delete data.reminder_time;
        if (data.frequency !== 'custom') {
            delete data.recurrence_days;
            delete data.recurrence_interval;
        }
        return data;
    };

    const toggleDay = (bit) => {
        setFormData((prev) => ({
            ...prev,
            recurrence_days: prev.recurrence_days ^ bit,
        }));
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Create New Habit">
            <form onSubmit={handleSubmit} className="space-y-5">
                <Input
                    label="Habit Name"
                    placeholder="e.g., Exercise, Read, Meditate"
                    value={formData.name}
                    onChange={handleChange('name')}
                    error={errors.name}
                />

                <TextArea
                    label="Description (Optional)"
                    placeholder="What's this habit about?"
                    value={formData.description}
                    onChange={handleChange('description')}
                    rows={3}
                />

                <div className="grid grid-cols-2 gap-4">
                    <Select label="Frequency" options={frequencyOptions} value={formData.frequency} onChange={handleChange('frequency')} />

                    <Input
                        label="Target Count"
                        type="number"
                        min="1"
                        value={formData.target_count}
                        onChange={handleChange('target_count')}
                        error={errors.target_count}
                        helperText="How many times per period?"
                    />
                </div>

                <Input
                    label="Daily Reminder (Optional)"
                    type="time"
                    value={formData.reminder_time}
                    onChange={handleChange('reminder_time')}
                    helperText="We'll send you a notification at this time"
                />

                <div className="flex justify-end gap-3 pt-4 border-t border-base-200 mt-6">
                    <Button type="button" variant="secondary" onClick={onClose}>
                        Cancel
                    </Button>
                    <Button type="submit" variant="primary" isLoading={isLoading}>
                        Create Habit
                    </Button>
                </div>
            </form>
        </Modal>
    );
}

export function EditHabitModal({ isOpen, onClose, habit }) {
    const { updateHabit, isLoading } = useHabitsStore();
    const { addToast } = useUIStore();
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        frequency: 'daily',
        target_count: 1,
        reminder_time: '',
    });
    const [errors, setErrors] = useState({});

    useEffect(() => {
        if (habit) {
            setFormData({
                name: habit.name || '',
                description: habit.description || '',
                frequency: habit.frequency || 'daily',
                target_count: habit.target_count || 1,
                reminder_time: habit.reminder_time || '',
            });
        }
    }, [habit]);

    const validate = () => {
        const newErrors = {};
        if (!formData.name.trim()) newErrors.name = 'Name is required';
        if (formData.target_count < 1) newErrors.target_count = 'Target must be at least 1';
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!validate()) return;

        const submitData = { ...formData };
        if (!submitData.reminder_time) delete submitData.reminder_time;

        const result = await updateHabit(habit.id, {
            ...submitData,
            target_count: parseInt(formData.target_count, 10),
        });

        if (result.success) {
            addToast({ type: 'success', title: 'Habit Updated', message: 'Your habit has been updated successfully.' });
            onClose();
        } else {
            addToast({ type: 'error', title: 'Error', message: result.error });
        }
    };

    const handleChange = (field) => (e) => {
        setFormData({ ...formData, [field]: e.target.value });
        if (errors[field]) setErrors({ ...errors, [field]: null });
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Edit Habit">
            <form onSubmit={handleSubmit} className="space-y-5">
                <Input
                    label="Habit Name"
                    placeholder="e.g., Exercise, Read, Meditate"
                    value={formData.name}
                    onChange={handleChange('name')}
                    error={errors.name}
                />

                <TextArea
                    label="Description (Optional)"
                    placeholder="What's this habit about?"
                    value={formData.description}
                    onChange={handleChange('description')}
                    rows={3}
                />

                <div className="grid grid-cols-2 gap-4">
                    <Select label="Frequency" options={frequencyOptions} value={formData.frequency} onChange={handleChange('frequency')} />

                    <Input
                        label="Target Count"
                        type="number"
                        min="1"
                        value={formData.target_count}
                        onChange={handleChange('target_count')}
                        error={errors.target_count}
                    />
                </div>

                <Input
                    label="Daily Reminder (Optional)"
                    type="time"
                    value={formData.reminder_time}
                    onChange={handleChange('reminder_time')}
                    helperText="We'll send you a notification at this time"
                />

                <div className="flex justify-end gap-3 pt-4 border-t border-base-200 mt-6">
                    <Button type="button" variant="secondary" onClick={onClose}>
                        Cancel
                    </Button>
                    <Button type="submit" variant="primary" isLoading={isLoading}>
                        Save Changes
                    </Button>
                </div>
            </form>
        </Modal>
    );
}

export function LogHabitModal({ isOpen, onClose, habit }) {
    const { logHabit, isLoading } = useHabitsStore();
    const { addToast } = useUIStore();
    const [count, setCount] = useState(1);
    const [notes, setNotes] = useState('');
    const [logDate, setLogDate] = useState(new Date().toISOString().split('T')[0]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!habit) return;

        const result = await logHabit(habit.id, {
            count: parseInt(count, 10),
            note: notes.trim() || undefined,
            log_date: logDate,
        });

        if (result.success) {
            addToast({ type: 'success', title: 'Progress Logged', message: 'Your progress has been recorded.' });
            setCount(1);
            setNotes('');
            onClose();
        } else {
            addToast({ type: 'error', title: 'Error', message: result.error });
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={`Log Progress: ${habit?.name || ''}`}>
            <form onSubmit={handleSubmit} className="space-y-5">
                <Input
                    label="Count"
                    type="number"
                    min="1"
                    value={count}
                    onChange={(e) => setCount(e.target.value)}
                    helperText="How many times did you complete this?"
                />

                <Input label="Date" type="date" value={logDate} onChange={(e) => setLogDate(e.target.value)} max={new Date().toISOString().split('T')[0]} />

                <TextArea
                    label="Notes (Optional)"
                    placeholder="Any notes about this session?"
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    rows={3}
                />

                <div className="flex justify-end gap-3 pt-4 border-t border-base-200 mt-6">
                    <Button type="button" variant="secondary" onClick={onClose}>
                        Cancel
                    </Button>
                    <Button type="submit" variant="success" isLoading={isLoading}>
                        Log Progress
                    </Button>
                </div>
            </form>
        </Modal>
    );
}

export function EditHabitLogModal({ isOpen, onClose, log }) {
    const { updateHabitLog, isLoading } = useHabitsStore();
    const { addToast } = useUIStore();
    const [count, setCount] = useState(1);
    const [notes, setNotes] = useState('');
    const [logDate, setLogDate] = useState('');

    useEffect(() => {
        if (log) {
            setCount(log.count || 1);
            setNotes(log.note || '');
            if (log.log_date) {
                // Ensure correct date format YYYY-MM-DD
                const date = new Date(log.log_date);
                setLogDate(date.toISOString().split('T')[0]);
            }
        }
    }, [log]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!log) return;

        const result = await updateHabitLog(log.id, {
            count: parseInt(count, 10),
            note: notes.trim() || undefined,
            log_date: logDate ? new Date(logDate).toISOString() : undefined,
        });

        if (result.success) {
            addToast({ type: 'success', title: 'Log Updated', message: 'Your progress log has been updated.' });
            onClose();
        } else {
            addToast({ type: 'error', title: 'Error', message: result.error });
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Edit Log">
            <form onSubmit={handleSubmit} className="space-y-5">
                <Input label="Count" type="number" min="1" value={count} onChange={(e) => setCount(e.target.value)} helperText="Update count" />

                <Input label="Date" type="date" value={logDate} onChange={(e) => setLogDate(e.target.value)} max={new Date().toISOString().split('T')[0]} />

                <TextArea label="Notes (Optional)" placeholder="Update notes..." value={notes} onChange={(e) => setNotes(e.target.value)} rows={3} />

                <div className="flex justify-end gap-3 pt-4 border-t border-base-200 mt-6">
                    <Button type="button" variant="secondary" onClick={onClose}>
                        Cancel
                    </Button>
                    <Button type="submit" variant="primary" isLoading={isLoading}>
                        Save Changes
                    </Button>
                </div>
            </form>
        </Modal>
    );
}
