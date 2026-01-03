import { useEffect, useState } from 'react';
import { format } from 'date-fns';
import { Target, Flame, TrendingUp, Calendar, Plus, CheckCircle2, Lightbulb } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Header } from '../../components/layout/Sidebar';
import { Card, StatsCard } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { PageLoader } from '../../components/ui/Loading';
import { HabitCardCompact, CreateHabitModal, LogHabitModal } from '../../components/habits';
import { useHabitsStore } from '../../stores/habitsStore';

export function DashboardPage() {
  const { t } = useTranslation();
  const { habits, dashboard, isLoading, fetchHabits, fetchDashboard } = useHabitsStore();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState(null);

  useEffect(() => {
    fetchHabits({ page: 1, per_page: 20, active: true });
    fetchDashboard();
  }, [fetchHabits, fetchDashboard]);

  const handleOpenLogModal = (habit) => {
    setSelectedHabit(habit);
    setIsLogModalOpen(true);
  };

  const activeHabits = habits?.filter((h) => h.is_active) || [];
  const today = format(new Date(), 'EEEE, MMMM d');

  if (isLoading && !habits.length) {
    return <PageLoader />;
  }

  return (
    <div className="space-y-8">
      <Header
        title={t('dashboard.title')}
        subtitle={today}
        actions={
          <Button variant="primary" onClick={() => setIsCreateModalOpen(true)}>
            <Plus size={16} />
            {t('habits.createHabit')}
          </Button>
        }
      />

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title={t('dashboard.activeHabits')}
          value={dashboard?.active_habits_count || activeHabits.length}
          subtitle={t('habits.title')}
          icon={Target}
        />
        <StatsCard title={t('dashboard.completedToday')} value={dashboard?.total_logs_today || 0} subtitle={t('common.success')} icon={CheckCircle2} />
        <StatsCard title={t('dashboard.currentStreak')} value={dashboard?.current_streak || 0} subtitle={t('dashboard.days')} icon={Flame} />
        <StatsCard
          title={t('dashboard.weeklyProgress')}
          value={`${dashboard?.weekly_completion || 0}%`}
          subtitle={t('habits.stats.completionRate')}
          icon={TrendingUp}
        />
      </div>

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Today's Habits (2/3 width) */}
        <div className="lg:col-span-2">
          <Card className="h-full">
            <div className="flex items-center gap-4 mb-6">
              <div className="p-2.5 bg-primary/5 rounded-lg">
                <Calendar size={20} className="text-primary" />
              </div>
              <div>
                <h3 className="text-base font-semibold text-base-content">{t('dashboard.todaysHabits')}</h3>
                <p className="text-xs text-base-content/50">
                  {activeHabits.length} {t('dashboard.activeToTrack')}
                </p>
              </div>
            </div>

            {activeHabits.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-16 text-center">
                <div className="p-5 bg-base-200 rounded-full mb-4">
                  <Target size={32} className="text-base-content/30" />
                </div>
                <h3 className="text-base font-semibold text-base-content mb-1">{t('dashboard.noHabitsYet')}</h3>
                <p className="text-sm text-base-content/50 mb-6 max-w-xs">{t('dashboard.createToStart')}</p>
                <Button variant="primary" onClick={() => setIsCreateModalOpen(true)}>
                  <Plus size={16} />
                  {t('dashboard.createFirst')}
                </Button>
              </div>
            ) : (
              <div className="space-y-1 -mx-6 -mb-6">
                {activeHabits.slice(0, 5).map((habit) => (
                  <HabitCardCompact key={habit.id} habit={habit} onLog={handleOpenLogModal} />
                ))}
                {activeHabits.length > 5 && (
                  <div className="px-6 py-3 text-center border-t border-base-200">
                    <button className="text-sm text-primary hover:text-primary/80 font-medium transition-colors">
                      +{activeHabits.length - 5} {t('dashboard.moreHabits')}
                    </button>
                  </div>
                )}
              </div>
            )}
          </Card>
        </div>

        {/* Sidebar Cards (1/3 width) */}
        <div className="flex flex-col gap-6">
          {/* Streak Card */}
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2 bg-warning/10 rounded-lg">
                <Flame size={18} className="text-warning" />
              </div>
              <div>
                <h3 className="text-sm font-semibold text-base-content">{t('dashboard.streakCard.title')}</h3>
                <p className="text-xs text-base-content/50">{t('dashboard.streakCard.subtitle')}</p>
              </div>
            </div>
            <div className="flex items-baseline gap-2 mb-3">
              <span className="text-3xl font-semibold text-base-content">{dashboard?.current_streak || 0}</span>
              <span className="text-sm text-base-content/50">{t('dashboard.days')}</span>
            </div>
            <div className="h-2 bg-base-200 rounded-full overflow-hidden">
              <div
                className="h-full bg-warning rounded-full transition-all duration-500"
                style={{ width: `${Math.min((dashboard?.current_streak || 0) * 10, 100)}%` }}
              />
            </div>
            <p className="text-xs text-base-content/50 mt-2 text-right">
              {t('dashboard.streakCard.nextMilestone')}: 10 {t('dashboard.days')}
            </p>
          </Card>

          {/* Quick Tips */}
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-2 bg-info/10 rounded-lg">
                <Lightbulb size={18} className="text-info" />
              </div>
              <h3 className="text-sm font-semibold text-base-content">{t('dashboard.tips.title')}</h3>
            </div>
            <ul className="space-y-4">
              {[
                { text: t('dashboard.tips.tip1'), icon: '🎯' },
                { text: t('dashboard.tips.tip2'), icon: '⏰' },
                { text: t('dashboard.tips.tip3'), icon: '🔗' },
              ].map((tip, index) => (
                <li key={index} className="flex gap-3 text-sm">
                  <span className="text-base select-none">{tip.icon}</span>
                  <span className="text-base-content/70 leading-relaxed">{tip.text}</span>
                </li>
              ))}
            </ul>
          </Card>
        </div>
      </div>

      {/* Modals */}
      <CreateHabitModal isOpen={isCreateModalOpen} onClose={() => setIsCreateModalOpen(false)} />
      <LogHabitModal isOpen={isLogModalOpen} onClose={() => setIsLogModalOpen(false)} habit={selectedHabit} />
    </div>
  );
}
