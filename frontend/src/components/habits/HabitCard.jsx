import { Link } from 'react-router-dom';
import { MoreVertical, Play, Pause, Trash2, Edit3, CheckCircle2, Flame } from 'lucide-react';
import { Menu, Transition } from '@headlessui/react';
import { Fragment } from 'react';
import { FrequencyBadge, StatusBadge } from '../ui/Badge';
import { useHabitsStore } from '../../stores/habitsStore';
import { useUIStore } from '../../stores/uiStore';

export function HabitCard({ habit, onLog, onEdit, onDelete }) {
  const { activateHabit, deactivateHabit } = useHabitsStore();
  const { addToast } = useUIStore();

  const handleToggleActive = async (e) => {
    e.preventDefault();
    e.stopPropagation();
    const action = habit.is_active ? deactivateHabit : activateHabit;
    const result = await action(habit.id);
    if (result.success) {
      addToast({ type: 'success', message: `Habit ${habit.is_active ? 'deactivated' : 'activated'} successfully` });
    } else {
      addToast({ type: 'error', message: result.error });
    }
  };

  const handleLog = (e) => {
    e.preventDefault();
    e.stopPropagation();
    onLog?.(habit);
  };

  return (
    <Link to={`/habits/${habit.id}`} className="block group">
      <div className="bg-base-100 border border-base-300 rounded-lg p-5 h-full transition-all duration-200 hover:shadow-md hover:border-base-300">
        <div className="flex justify-between items-start mb-3">
          <div className="flex-1 min-w-0 mr-2">
            <div className="flex items-center gap-2 mb-1.5">
              <h3 className="text-base font-semibold text-base-content truncate" title={habit.name}>
                {habit.name}
              </h3>
              <StatusBadge isActive={habit.is_active} />
              {habit.vacation_mode && <span className="badge badge-accent badge-sm uppercase text-[10px]">Vacation</span>}
            </div>
            <p className="text-sm text-base-content/50 line-clamp-2 leading-relaxed min-h-[2.5em]">{habit.description || 'No description added'}</p>
          </div>

          {/* Action Menu */}
          <Menu as="div" className="relative ml-2 shrink-0">
            <Menu.Button
              onClick={(e) => e.preventDefault()}
              className="p-1.5 rounded-md text-base-content/50 hover:text-base-content hover:bg-base-200 transition-colors"
            >
              <MoreVertical size={16} />
            </Menu.Button>
            <Transition
              as={Fragment}
              enter="transition ease-out duration-100"
              enterFrom="transform opacity-0 scale-95"
              enterTo="transform opacity-100 scale-100"
              leave="transition ease-in duration-75"
              leaveFrom="transform opacity-100 scale-100"
              leaveTo="transform opacity-0 scale-95"
            >
              <Menu.Items className="absolute right-0 mt-1 w-44 origin-top-right bg-base-100 border border-base-300 rounded-lg shadow-lg focus:outline-none z-10 py-1">
                <Menu.Item>
                  {({ active }) => (
                    <button
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        onEdit?.(habit);
                      }}
                      className={`flex w-full items-center px-3 py-2 text-sm gap-2.5 ${active ? 'bg-base-200' : ''} text-base-content`}
                    >
                      <Edit3 size={14} />
                      Edit Habit
                    </button>
                  )}
                </Menu.Item>
                <Menu.Item>
                  {({ active }) => (
                    <button
                      onClick={handleToggleActive}
                      className={`flex w-full items-center px-3 py-2 text-sm gap-2.5 ${active ? 'bg-base-200' : ''} text-base-content`}
                    >
                      {habit.is_active ? <Pause size={14} /> : <Play size={14} />}
                      {habit.is_active ? 'Deactivate' : 'Activate'}
                    </button>
                  )}
                </Menu.Item>
                <div className="border-t border-base-200 my-1" />
                <Menu.Item>
                  {({ active }) => (
                    <button
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        onDelete?.(habit);
                      }}
                      className={`flex w-full items-center px-3 py-2 text-sm gap-2.5 text-error ${active ? 'bg-error/5' : ''}`}
                    >
                      <Trash2 size={14} />
                      Delete Habit
                    </button>
                  )}
                </Menu.Item>
              </Menu.Items>
            </Transition>
          </Menu>
        </div>

        <div className="flex items-center gap-3 mb-4">
          <FrequencyBadge frequency={habit.frequency} />
          <span className="text-xs text-base-content/50 font-medium">Target: {habit.target_count}x</span>
          {habit.current_streak > 0 && (
            <div className="flex items-center gap-1 text-orange-500" title="Current Streak">
              <Flame size={14} fill="currentColor" />
              <span className="text-xs font-bold">{habit.current_streak}</span>
            </div>
          )}
        </div>

        {habit.is_active && (
          <button
            className="
                            w-full flex items-center justify-center gap-2
                            py-2 text-sm font-medium
                            bg-primary/5 text-primary border border-primary/20 rounded-md
                            hover:bg-primary/10 transition-colors
                        "
            onClick={handleLog}
          >
            <CheckCircle2 size={16} />
            Log Progress
          </button>
        )}
      </div>
    </Link>
  );
}

export function HabitCardCompact({ habit, onLog }) {
  return (
    <div className="flex items-center justify-between px-6 py-4 border-b border-base-200 last:border-none hover:bg-base-50 transition-colors">
      <div className="flex items-center gap-4">
        <div className={`w-2 h-2 rounded-full shrink-0 ${habit.is_active ? 'bg-success' : 'bg-base-300'}`} />
        <div>
          <p className="text-sm font-medium text-base-content">{habit.name}</p>
          <div className="flex items-center gap-2 text-xs text-base-content/50 mt-0.5">
            <span className="capitalize">{habit.frequency}</span>
            <span>•</span>
            <span>Target: {habit.target_count}x</span>
          </div>
        </div>
      </div>

      <button className="p-2 rounded-md text-success hover:bg-success/10 transition-colors" onClick={() => onLog?.(habit)} title="Log Progress">
        <CheckCircle2 size={20} />
      </button>
    </div>
  );
}

export function HabitStats({ stats }) {
  const statItems = [
    { label: 'Total Logs', value: stats?.total_logs || 0 },
    { label: 'Current Streak', value: stats?.current_streak || 0, icon: Flame, highlight: true },
    { label: 'Longest Streak', value: stats?.longest_streak || 0 },
  ];

  return (
    <div className="grid grid-cols-3 gap-4">
      {statItems.map((item) => (
        <div key={item.label} className="bg-base-100 border border-base-300 rounded-lg p-4 text-center">
          <p className="text-xs font-medium text-base-content/50 uppercase tracking-wide mb-1">{item.label}</p>
          <div className={`flex items-center justify-center gap-1.5 ${item.highlight ? 'text-warning' : 'text-base-content'}`}>
            {item.icon && <item.icon size={18} />}
            <span className="text-xl font-semibold">{item.value}</span>
          </div>
        </div>
      ))}
    </div>
  );
}
