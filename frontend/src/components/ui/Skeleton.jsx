/**
 * Skeleton Loading Components
 * Beautiful shimmer placeholders for loading states
 */

// Base Skeleton with shimmer animation
export function Skeleton({ className = '', animate = true }) {
  return (
    <div
      className={`
                bg-base-300 rounded
                ${animate ? 'animate-pulse' : ''}
                ${className}
            `}
    />
  );
}

// Text line skeleton
export function SkeletonText({ lines = 1, className = '' }) {
  return (
    <div className={`space-y-2 ${className}`}>
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton key={i} className={`h-4 ${i === lines - 1 && lines > 1 ? 'w-3/4' : 'w-full'}`} />
      ))}
    </div>
  );
}

// Avatar skeleton
export function SkeletonAvatar({ size = 'md', className = '' }) {
  const sizes = {
    sm: 'w-8 h-8',
    md: 'w-10 h-10',
    lg: 'w-12 h-12',
    xl: 'w-16 h-16',
  };

  return <Skeleton className={`${sizes[size]} rounded-full ${className}`} />;
}

// Card skeleton
export function SkeletonCard({ className = '' }) {
  return (
    <div className={`bg-base-100 border border-base-300 rounded-lg p-4 ${className}`}>
      <div className="flex items-start gap-4">
        <SkeletonAvatar size="md" />
        <div className="flex-1 space-y-3">
          <Skeleton className="h-5 w-1/3" />
          <SkeletonText lines={2} />
        </div>
      </div>
    </div>
  );
}

// Habit card skeleton
export function SkeletonHabitCard() {
  return (
    <div className="bg-base-100 border border-base-300 rounded-lg p-4">
      <div className="flex items-center gap-4">
        <Skeleton className="w-10 h-10 rounded-lg" />
        <div className="flex-1">
          <Skeleton className="h-5 w-1/2 mb-2" />
          <Skeleton className="h-3 w-3/4" />
        </div>
        <Skeleton className="w-6 h-6 rounded" />
      </div>
      <div className="mt-4 flex gap-1">
        {Array.from({ length: 7 }).map((_, i) => (
          <Skeleton key={i} className="flex-1 h-8 rounded" />
        ))}
      </div>
    </div>
  );
}

// Stats card skeleton
export function SkeletonStatsCard() {
  return (
    <div className="bg-base-100 border border-base-300 rounded-lg p-5">
      <div className="flex items-center justify-between mb-3">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="w-8 h-8 rounded-lg" />
      </div>
      <Skeleton className="h-8 w-20 mb-2" />
      <Skeleton className="h-3 w-16" />
    </div>
  );
}

// Dashboard skeleton
export function SkeletonDashboard() {
  return (
    <div className="space-y-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <Skeleton className="h-7 w-48 mb-2" />
          <Skeleton className="h-4 w-32" />
        </div>
        <Skeleton className="h-10 w-32 rounded-lg" />
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <SkeletonStatsCard />
        <SkeletonStatsCard />
        <SkeletonStatsCard />
        <SkeletonStatsCard />
      </div>

      {/* Habits */}
      <div>
        <Skeleton className="h-6 w-32 mb-4" />
        <div className="grid gap-4">
          <SkeletonHabitCard />
          <SkeletonHabitCard />
          <SkeletonHabitCard />
        </div>
      </div>
    </div>
  );
}

// Habits list skeleton
export function SkeletonHabitsList() {
  return (
    <div className="space-y-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="flex justify-between items-center pb-6 border-b border-base-200">
        <div>
          <Skeleton className="h-6 w-24 mb-2" />
          <Skeleton className="h-4 w-40" />
        </div>
        <Skeleton className="h-10 w-28 rounded-lg" />
      </div>

      {/* Habits Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <SkeletonHabitCard key={i} />
        ))}
      </div>
    </div>
  );
}

// Table row skeleton
export function SkeletonTableRow({ columns = 4 }) {
  return (
    <tr>
      {Array.from({ length: columns }).map((_, i) => (
        <td key={i} className="py-4 px-4">
          <Skeleton className="h-4 w-full" />
        </td>
      ))}
    </tr>
  );
}

// Table skeleton
export function SkeletonTable({ rows = 5, columns = 4 }) {
  return (
    <div className="bg-base-100 border border-base-300 rounded-lg overflow-hidden">
      <table className="w-full">
        <thead>
          <tr className="bg-base-200">
            {Array.from({ length: columns }).map((_, i) => (
              <th key={i} className="py-3 px-4 text-left">
                <Skeleton className="h-4 w-20" />
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: rows }).map((_, i) => (
            <SkeletonTableRow key={i} columns={columns} />
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Analytics skeleton
export function SkeletonAnalytics() {
  return (
    <div className="space-y-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="pb-6 border-b border-base-200">
        <Skeleton className="h-6 w-24 mb-2" />
        <Skeleton className="h-4 w-48" />
      </div>

      {/* Charts */}
      <div className="grid gap-6 lg:grid-cols-2">
        <div className="bg-base-100 border border-base-300 rounded-lg p-5">
          <Skeleton className="h-5 w-32 mb-4" />
          <Skeleton className="h-64 w-full rounded-lg" />
        </div>
        <div className="bg-base-100 border border-base-300 rounded-lg p-5">
          <Skeleton className="h-5 w-32 mb-4" />
          <Skeleton className="h-64 w-full rounded-lg" />
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <SkeletonStatsCard />
        <SkeletonStatsCard />
        <SkeletonStatsCard />
        <SkeletonStatsCard />
      </div>
    </div>
  );
}

// Settings skeleton
export function SkeletonSettings() {
  return (
    <div className="space-y-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="pb-6 border-b border-base-200">
        <Skeleton className="h-6 w-24 mb-2" />
        <Skeleton className="h-4 w-48" />
      </div>

      {/* Settings Sections */}
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className="bg-base-100 border border-base-300 rounded-lg p-6">
          <div className="flex items-start gap-4">
            <Skeleton className="w-10 h-10 rounded-lg" />
            <div className="flex-1">
              <Skeleton className="h-5 w-40 mb-2" />
              <Skeleton className="h-4 w-64 mb-4" />
              <div className="grid gap-4 md:grid-cols-2">
                <Skeleton className="h-10 rounded-lg" />
                <Skeleton className="h-10 rounded-lg" />
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
