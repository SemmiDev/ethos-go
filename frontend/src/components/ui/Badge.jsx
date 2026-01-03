export function Badge({ children, variant = 'neutral', size = 'sm', className = '' }) {
    const variants = {
        gray: 'bg-base-200 text-base-content/70',
        neutral: 'bg-base-200 text-base-content/70',
        primary: 'bg-primary/10 text-primary',
        blue: 'bg-info/10 text-info',
        green: 'bg-success/10 text-success',
        success: 'bg-success/10 text-success',
        amber: 'bg-warning/10 text-warning',
        warning: 'bg-warning/10 text-warning',
        red: 'bg-error/10 text-error',
        danger: 'bg-error/10 text-error',
        error: 'bg-error/10 text-error',
        info: 'bg-info/10 text-info',
        teal: 'bg-accent/10 text-accent',
        accent: 'bg-accent/10 text-accent',
    };

    const sizes = {
        xs: 'px-1.5 py-0.5 text-[10px]',
        sm: 'px-2 py-0.5 text-xs',
        md: 'px-2.5 py-1 text-xs',
    };

    const variantClass = variants[variant] || variants.neutral;
    const sizeClass = sizes[size] || sizes.sm;

    return <span className={`inline-flex items-center font-medium rounded ${variantClass} ${sizeClass} ${className}`}>{children}</span>;
}

export function FrequencyBadge({ frequency }) {
    const config = {
        daily: { label: 'Daily', variant: 'primary' },
        weekly: { label: 'Weekly', variant: 'info' },
        monthly: { label: 'Monthly', variant: 'neutral' },
    };

    const { label, variant } = config[frequency] || { label: frequency, variant: 'neutral' };

    return <Badge variant={variant}>{label}</Badge>;
}

export function StatusBadge({ isActive }) {
    return (
        <Badge variant={isActive ? 'success' : 'neutral'} className="gap-1.5">
            <span className={`h-1.5 w-1.5 rounded-full ${isActive ? 'bg-success' : 'bg-base-content/40'}`} />
            {isActive ? 'Active' : 'Inactive'}
        </Badge>
    );
}

export function CountBadge({ count, max }) {
    const percentage = max > 0 ? (count / max) * 100 : 0;
    let variant = 'neutral';
    if (percentage >= 100) variant = 'success';
    else if (percentage >= 50) variant = 'warning';
    else if (percentage > 0) variant = 'info';

    return (
        <Badge variant={variant}>
            {count}/{max}
        </Badge>
    );
}
