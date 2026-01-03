export function Spinner({ size = 'md', className = '' }) {
    const sizes = {
        xs: 'h-3 w-3',
        sm: 'h-4 w-4',
        md: 'h-5 w-5',
        lg: 'h-6 w-6',
        xl: 'h-8 w-8',
    };

    const sizeClass = sizes[size] || sizes.md;

    return (
        <svg className={`animate-spin ${sizeClass} ${className}`} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
        </svg>
    );
}

export function LoadingScreen({ message = 'Loading...' }) {
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-base-100">
            <div className="text-center">
                <Spinner size="xl" className="text-primary mx-auto mb-4" />
                <p className="font-medium text-base-content">{message}</p>
                <p className="text-sm text-base-content/50 mt-1">Please wait a moment</p>
            </div>
        </div>
    );
}

export function PageLoader() {
    return (
        <div className="flex flex-col items-center justify-center min-h-[50vh] gap-4">
            <Spinner size="lg" className="text-primary" />
            <p className="text-sm text-base-content/50">Loading content...</p>
        </div>
    );
}

export function Skeleton({ className = '', style = {} }) {
    return <div className={`bg-base-200 rounded animate-pulse ${className}`} style={{ minHeight: 20, ...style }} />;
}

export function CardSkeleton() {
    return (
        <div className="bg-base-100 border border-base-300 rounded-lg p-6">
            <div className="flex justify-between items-start mb-4">
                <div className="w-full">
                    <Skeleton className="w-1/3 h-5 mb-2" />
                    <Skeleton className="w-2/3 h-4" />
                </div>
                <Skeleton className="w-16 h-6 rounded" />
            </div>
            <Skeleton className="w-full h-4 mb-4" />
            <div className="flex gap-3">
                <Skeleton className="w-20 h-9 rounded-md" />
                <Skeleton className="w-20 h-9 rounded-md" />
            </div>
        </div>
    );
}

export function StatsSkeleton() {
    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {[...Array(4)].map((_, i) => (
                <div key={i} className="bg-base-100 border border-base-300 rounded-lg p-5">
                    <div className="flex items-start justify-between">
                        <div className="flex-1">
                            <Skeleton className="w-20 h-4 mb-2" />
                            <Skeleton className="w-16 h-7" />
                        </div>
                        <Skeleton className="w-11 h-11 rounded-lg" />
                    </div>
                </div>
            ))}
        </div>
    );
}
