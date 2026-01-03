import { useUIStore } from '../../stores/uiStore';
import { X, CheckCircle, XCircle, Info, AlertTriangle } from 'lucide-react';

const icons = {
    success: CheckCircle,
    error: XCircle,
    info: Info,
    warning: AlertTriangle,
};

const variants = {
    success: {
        bg: 'bg-success/5',
        border: 'border-success/20',
        text: 'text-success',
        icon: 'text-success',
    },
    error: {
        bg: 'bg-error/5',
        border: 'border-error/20',
        text: 'text-error',
        icon: 'text-error',
    },
    info: {
        bg: 'bg-info/5',
        border: 'border-info/20',
        text: 'text-info',
        icon: 'text-info',
    },
    warning: {
        bg: 'bg-warning/5',
        border: 'border-warning/20',
        text: 'text-warning',
        icon: 'text-warning',
    },
};

export function ToastContainer() {
    const { toasts, removeToast } = useUIStore();

    if (!toasts.length) return null;

    return (
        <div className="fixed top-4 right-4 z-[100] flex flex-col gap-3 max-w-sm w-full">
            {toasts.map((toast) => (
                <Toast key={toast.id} toast={toast} onClose={() => removeToast(toast.id)} />
            ))}
        </div>
    );
}

function Toast({ toast, onClose }) {
    const Icon = icons[toast.type] || icons.info;
    const variant = variants[toast.type] || variants.info;

    return (
        <div
            className={`
                flex items-start gap-3 p-4 rounded-lg border shadow-lg
                ${variant.bg} ${variant.border}
                animate-in slide-in-from-right fade-in duration-300
            `}
        >
            <Icon className={`h-5 w-5 shrink-0 mt-0.5 ${variant.icon}`} />
            <div className="flex-1 min-w-0">
                {toast.title && <h4 className={`text-sm font-semibold ${variant.text}`}>{toast.title}</h4>}
                <p className="text-sm text-base-content/70 mt-0.5">{toast.message}</p>
            </div>
            <button className="p-1 rounded-md hover:bg-base-200/50 transition-colors shrink-0" onClick={onClose}>
                <X size={14} className="text-base-content/50" />
            </button>
        </div>
    );
}
