import { forwardRef } from 'react';

export const Button = forwardRef(({ children, variant = 'primary', size = 'md', isLoading = false, disabled = false, className = '', ...props }, ref) => {
    const variants = {
        primary: 'bg-primary text-primary-content hover:bg-primary/90 border-primary',
        secondary: 'bg-base-100 text-base-content hover:bg-base-200 border border-base-300',
        outline: 'bg-transparent text-primary border border-primary hover:bg-primary/5',
        ghost: 'bg-transparent text-base-content hover:bg-base-200 border-transparent',
        danger: 'bg-error text-error-content hover:bg-error/90 border-error',
        success: 'bg-success text-success-content hover:bg-success/90 border-success',
        info: 'bg-info text-info-content hover:bg-info/90 border-info',
        warning: 'bg-warning text-warning-content hover:bg-warning/90 border-warning',
        neutral: 'bg-neutral text-neutral-content hover:bg-neutral/90 border-neutral',
        accent: 'bg-accent text-accent-content hover:bg-accent/90 border-accent',
    };

    const sizes = {
        xs: 'px-2.5 py-1 text-xs',
        sm: 'px-3 py-1.5 text-sm',
        md: 'px-4 py-2 text-sm',
        lg: 'px-5 py-2.5 text-base',
        icon: 'p-2',
    };

    const variantClass = variants[variant] || variants.primary;
    const sizeClass = sizes[size] || sizes.md;

    return (
        <button
            ref={ref}
            disabled={disabled || isLoading}
            className={`
                inline-flex items-center justify-center gap-2
                font-medium rounded-md
                transition-all duration-150 ease-in-out
                focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/30 focus-visible:ring-offset-2
                disabled:opacity-50 disabled:cursor-not-allowed
                ${variantClass}
                ${sizeClass}
                ${className}
            `}
            {...props}
        >
            {isLoading && (
                <svg className="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                </svg>
            )}
            {children}
        </button>
    );
});

Button.displayName = 'Button';
