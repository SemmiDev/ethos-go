import { format } from 'date-fns';
import { Calendar, FileText, Edit2, Trash2 } from 'lucide-react';
import { Badge } from '../ui/Badge';

export function HabitLogList({ logs = [], onEdit, onDelete }) {
    if (!logs.length) {
        return (
            <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="p-4 bg-base-200 rounded-full mb-4">
                    <FileText size={28} className="text-base-content/30" />
                </div>
                <h3 className="font-semibold text-base text-base-content mb-1">No logs yet</h3>
                <p className="text-sm text-base-content/50">Start logging your progress to see your history here.</p>
            </div>
        );
    }

    return (
        <div className="divide-y divide-base-200 -mx-6">
            {logs.map((log) => (
                <div key={log.id} className="flex items-start gap-4 px-6 py-4 hover:bg-base-50 transition-colors group">
                    <div className="p-2.5 bg-primary/5 rounded-lg shrink-0">
                        <Calendar size={16} className="text-primary" />
                    </div>
                    <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                            <p className="text-sm font-medium text-base-content">{format(new Date(log.log_date), 'MMM d, yyyy')}</p>
                            <Badge variant="primary" size="xs">
                                ×{log.count}
                            </Badge>
                        </div>
                        {log.created_at && <p className="text-xs text-base-content/50">Logged at {format(new Date(log.created_at), 'h:mm a')}</p>}
                        {log.note && <p className="text-sm mt-2 p-2.5 bg-base-100 border border-base-200 rounded-md text-base-content/70">{log.note}</p>}
                    </div>
                    <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                            className="p-1.5 rounded-md text-base-content/50 hover:text-base-content hover:bg-base-200 transition-colors"
                            onClick={() => onEdit && onEdit(log)}
                            title="Edit Log"
                        >
                            <Edit2 size={14} />
                        </button>
                        <button
                            className="p-1.5 rounded-md text-error/70 hover:text-error hover:bg-error/5 transition-colors"
                            onClick={() => onDelete && onDelete(log)}
                            title="Delete Log"
                        >
                            <Trash2 size={14} />
                        </button>
                    </div>
                </div>
            ))}
        </div>
    );
}
