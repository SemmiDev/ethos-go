export function Card({ children, className = '', hover = false, ...props }) {
    return (
        <div
            className={`
                bg-base-100 border border-base-300 rounded-lg
                ${hover ? 'transition-shadow duration-200 hover:shadow-md' : ''}
                ${className}
            `}
            {...props}
        >
            <div className="p-6">{children}</div>
        </div>
    );
}

export function CardHeader({ children, className = '' }) {
    return <div className={`flex flex-col space-y-1 mb-5 ${className}`}>{children}</div>;
}

export function CardTitle({ children, className = '' }) {
    return <h3 className={`text-lg font-semibold text-base-content tracking-tight ${className}`}>{children}</h3>;
}

export function CardDescription({ children, className = '' }) {
    return <p className={`text-sm text-base-content/60 ${className}`}>{children}</p>;
}

export function CardContent({ children, className = '' }) {
    return <div className={className}>{children}</div>;
}

export function CardFooter({ children, className = '' }) {
    return <div className={`flex items-center pt-5 mt-5 border-t border-base-200 ${className}`}>{children}</div>;
}

export function StatsCard({ title, value, subtitle, icon: Icon, trend }) {
    return (
        <div className="bg-base-100 border border-base-300 rounded-lg p-5 transition-shadow duration-200 hover:shadow-sm">
            <div className="flex items-start justify-between">
                <div className="flex-1">
                    <p className="text-sm font-medium text-base-content/60 mb-1">{title}</p>
                    <p className="text-2xl font-semibold text-base-content tracking-tight">{value}</p>
                    {subtitle && <p className="text-xs text-base-content/50 mt-1.5">{subtitle}</p>}
                    {trend && (
                        <div
                            className={`flex items-center gap-1 mt-2 text-xs font-medium ${
                                trend > 0 ? 'text-success' : trend < 0 ? 'text-error' : 'text-base-content/50'
                            }`}
                        >
                            {trend > 0 ? '↑' : trend < 0 ? '↓' : '→'} {Math.abs(trend)}%
                        </div>
                    )}
                </div>
                {Icon && (
                    <div className="p-2.5 bg-primary/5 rounded-lg">
                        <Icon size={22} className="text-primary" />
                    </div>
                )}
            </div>
        </div>
    );
}

export function FeatureCard({ icon: Icon, title, description }) {
    return (
        <div className="bg-base-100 border border-base-300 rounded-lg p-5 transition-shadow duration-200 hover:shadow-md">
            <div className="mb-4 inline-flex items-center justify-center w-11 h-11 rounded-lg bg-primary/5">
                <Icon size={22} className="text-primary" />
            </div>
            <h3 className="text-base font-semibold text-base-content mb-2">{title}</h3>
            <p className="text-sm text-base-content/60 leading-relaxed">{description}</p>
        </div>
    );
}
