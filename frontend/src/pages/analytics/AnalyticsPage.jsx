import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { BarChart3, TrendingUp, Target, Calendar, Trophy } from 'lucide-react';
import { Header } from '../../components/layout/Sidebar';
import { Card, StatsCard } from '../../components/ui/Card';
import { Badge } from '../../components/ui/Badge';
import { PageLoader } from '../../components/ui/Loading';
import { useHabitsStore } from '../../stores/habitsStore';

export function AnalyticsPage() {
  const { t } = useTranslation();
  const { habits, dashboard, weeklyAnalytics, isLoading, fetchHabits, fetchDashboard, fetchWeeklyAnalytics } = useHabitsStore();

  useEffect(() => {
    fetchHabits();
    fetchDashboard();
    fetchWeeklyAnalytics();
  }, [fetchHabits, fetchDashboard, fetchWeeklyAnalytics]);

  if (isLoading && !habits.length) {
    return <PageLoader />;
  }

  // Use real data from API or fallback to zeros
  const weeklyData = weeklyAnalytics?.days?.map((day) => day.completion_percentage) || [0, 0, 0, 0, 0, 0, 0];
  const weeklyDaysLabels = weeklyAnalytics?.days?.map((day) => day.day_name) || ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

  return (
    <div className="space-y-6">
      <Header title={t('analytics.title')} subtitle={t('analytics.subtitle')} />

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard title={t('analytics.stats.totalHabits')} value={habits.length} subtitle={t('analytics.stats.createdHabits')} icon={Target} />
        <StatsCard
          title={t('analytics.stats.weeklyCompletion')}
          value={`${dashboard?.weekly_completion || weeklyAnalytics?.average_completion || 0}%`}
          subtitle={t('analytics.stats.thisWeek')}
          icon={TrendingUp}
        />
        <StatsCard title={t('analytics.stats.bestStreak')} value={dashboard?.longest_streak || 0} subtitle={t('analytics.stats.daysRecord')} icon={Trophy} />
        <StatsCard title={t('analytics.stats.totalLogs')} value={dashboard?.total_logs || 0} subtitle={t('analytics.stats.allTime')} icon={Calendar} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Weekly Progress Chart */}
        <div className="lg:col-span-2">
          <Card className="h-full">
            <div className="flex items-center gap-3 mb-6">
              <div className="p-2.5 bg-primary/5 rounded-lg">
                <BarChart3 size={20} className="text-primary" />
              </div>
              <div>
                <h3 className="text-base font-semibold text-base-content">{t('analytics.weeklyProgress')}</h3>
                <p className="text-xs text-base-content/50">{t('analytics.weeklyProgressDesc')}</p>
              </div>
            </div>

            {/* Simple Bar Chart */}
            <div className="flex items-end justify-between h-48 pt-4">
              {weeklyDaysLabels.map((day, index) => (
                <div key={day} className="flex flex-col items-center flex-1 group">
                  <div className="relative w-8">
                    <div
                      className="w-full bg-primary/20 group-hover:bg-primary transition-colors duration-200 rounded-t"
                      style={{ height: `${Math.max(weeklyData[index] || 0, 4)}%` }}
                    />
                    {/* Tooltip */}
                    <div className="absolute -top-8 left-1/2 -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-base-content text-base-100 text-xs px-2 py-1 rounded whitespace-nowrap">
                      {weeklyData[index] || 0}%
                    </div>
                  </div>
                  <span className="text-xs text-base-content/50 mt-2">{day}</span>
                </div>
              ))}
            </div>

            {/* Average line */}
            <div className="flex items-center justify-center gap-2 mt-4 pt-4 border-t border-base-200">
              <div className="h-2 w-2 rounded-full bg-primary/50" />
              <span className="text-xs text-base-content/60">
                {t('analytics.averageCompletion')}: {weeklyAnalytics?.average_completion || 0}%
              </span>
            </div>
          </Card>
        </div>

        {/* Habit Distribution */}
        <div>
          <Card className="h-full">
            <div className="flex items-center gap-3 mb-6">
              <div className="p-2.5 bg-secondary/5 rounded-lg">
                <Target size={20} className="text-secondary" />
              </div>
              <h3 className="text-base font-semibold text-base-content">{t('analytics.habitDistribution')}</h3>
            </div>

            <div className="space-y-5">
              {[
                { label: t('analytics.frequency.daily'), count: habits.filter((h) => h.frequency === 'daily').length, color: 'bg-primary' },
                { label: t('analytics.frequency.weekly'), count: habits.filter((h) => h.frequency === 'weekly').length, color: 'bg-secondary' },
                { label: t('analytics.frequency.monthly'), count: habits.filter((h) => h.frequency === 'monthly').length, color: 'bg-warning' },
              ].map((item) => (
                <div key={item.label}>
                  <div className="flex justify-between mb-2">
                    <span className="text-sm font-medium text-base-content">{item.label}</span>
                    <span className="text-sm text-base-content/50">{item.count}</span>
                  </div>
                  <div className="h-2 bg-base-200 rounded-full overflow-hidden">
                    <div
                      className={`h-full ${item.color} rounded-full transition-all duration-500`}
                      style={{ width: habits.length > 0 ? `${(item.count / habits.length) * 100}%` : '0%' }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </div>

      {/* Achievements */}
      <Card>
        <div className="flex items-center gap-3 mb-6">
          <div className="p-2.5 bg-warning/10 rounded-lg">
            <Trophy size={20} className="text-warning" />
          </div>
          <div>
            <h3 className="text-base font-semibold text-base-content">{t('analytics.achievements')}</h3>
            <p className="text-xs text-base-content/50">{t('analytics.achievementsDesc')}</p>
          </div>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {[
            { icon: '🎯', titleKey: 'firstHabit', unlocked: habits.length > 0 },
            { icon: '🔥', titleKey: 'sevenDayStreak', unlocked: (dashboard?.longest_streak || 0) >= 7 },
            { icon: '🏆', titleKey: 'thirtyDayMaster', unlocked: (dashboard?.longest_streak || 0) >= 30 },
            { icon: '⭐', titleKey: 'hundredLogs', unlocked: (dashboard?.total_logs || 0) >= 100 },
          ].map((achievement, index) => (
            <div
              key={index}
              className={`
                                flex flex-col items-center text-center p-5 rounded-lg border
                                ${achievement.unlocked ? 'border-base-300 bg-base-100' : 'border-base-200 bg-base-100 opacity-40 grayscale'}
                            `}
            >
              <div className="text-3xl mb-3">{achievement.icon}</div>
              <p className="font-semibold text-sm text-base-content mb-1">{t(`analytics.achievementsList.${achievement.titleKey}.title`)}</p>
              <p className="text-xs text-base-content/50 mb-3">{t(`analytics.achievementsList.${achievement.titleKey}.desc`)}</p>
              {achievement.unlocked && (
                <Badge variant="success" size="xs">
                  {t('analytics.achievementsList.unlocked')}
                </Badge>
              )}
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
