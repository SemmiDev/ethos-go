import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Edit3, Trash2, CheckCircle2, Activity, Clock, TrendingUp, Lightbulb } from 'lucide-react';
import { Header } from '../../components/layout/Sidebar';
import { Button } from '../../components/ui/Button';
import { Card } from '../../components/ui/Card';
import { Pagination } from '../../components/ui/Pagination';
import { FrequencyBadge, StatusBadge } from '../../components/ui/Badge';
import { PageLoader } from '../../components/ui/Loading';
import { ConfirmModal } from '../../components/ui/Modal';
import { HabitStats, HabitLogList, EditHabitModal, LogHabitModal, EditHabitLogModal } from '../../components/habits';
import { useHabitsStore } from '../../stores/habitsStore';
import { useUIStore } from '../../stores/uiStore';

export function HabitDetailPage() {
    const { t } = useTranslation();

    const params = useParams();
    const id = params.id || params.habitId;
    const navigate = useNavigate();
    const {
        selectedHabit,
        habitLogs,
        habitStats,
        logsPaging,
        isLoading,
        error,
        fetchHabit,
        fetchHabitLogs,
        fetchHabitStats,
        deleteHabit,
        deleteHabitLog,
        setLogsPage,
    } = useHabitsStore();

    const { addToast } = useUIStore();

    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isLogModalOpen, setIsLogModalOpen] = useState(false);
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    const [selectedLog, setSelectedLog] = useState(null);
    const [isEditLogModalOpen, setIsEditLogModalOpen] = useState(false);
    const [isDeleteLogModalOpen, setIsDeleteLogModalOpen] = useState(false);

    useEffect(() => {
        if (id) {
            fetchHabit(id);
            fetchHabitLogs(id, { page: 1 });
            fetchHabitStats(id);
        }
    }, [id, fetchHabit, fetchHabitLogs, fetchHabitStats]);

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[50vh] text-center p-4">
                <div className="p-4 bg-error/10 rounded-full mb-4">
                    <Trash2 size={28} className="text-error" />
                </div>
                <h3 className="text-lg font-semibold text-base-content mb-2">{t('common.error')}</h3>
                <p className="text-sm text-base-content/50 mb-6">{error}</p>
                <Button variant="secondary" onClick={() => navigate('/habits')}>
                    <ArrowLeft size={16} />
                    {t('common.back')}
                </Button>
            </div>
        );
    }

    const handleDelete = async () => {
        setIsDeleting(true);
        const result = await deleteHabit(id);
        setIsDeleting(false);

        if (result.success) {
            addToast({ type: 'success', title: t('common.success'), message: t('toast.habitDeleted') });
            navigate('/habits');
        } else {
            addToast({ type: 'error', title: t('common.error'), message: result.error });
        }
    };

    const handleEditLog = (log) => {
        setSelectedLog(log);
        setIsEditLogModalOpen(true);
    };

    const handleDeleteLog = (log) => {
        setSelectedLog(log);
        setIsDeleteLogModalOpen(true);
    };

    const confirmDeleteLog = async () => {
        if (!selectedLog) return;
        const result = await deleteHabitLog(selectedLog.id);
        if (result.success) {
            addToast({ type: 'success', title: t('common.success'), message: t('toast.logDeleted') });
            setIsDeleteLogModalOpen(false);
            setSelectedLog(null);
        } else {
            addToast({ type: 'error', title: t('common.error'), message: result.error });
        }
    };

    const handlePageChange = (page) => {
        setLogsPage(id, page);
    };

    if (isLoading || !selectedHabit) {
        return <PageLoader />;
    }

    return (
        <div className="space-y-6">
            <Header
                title={selectedHabit.name}
                actions={
                    <div className="flex gap-2">
                        <Button variant="ghost" size="sm" onClick={() => navigate('/habits')}>
                            <ArrowLeft size={16} />
                            {t('common.back')}
                        </Button>
                        <Button variant="secondary" size="sm" onClick={() => setIsEditModalOpen(true)}>
                            <Edit3 size={16} />
                            {t('common.edit')}
                        </Button>
                        <Button variant="danger" size="sm" onClick={() => setIsDeleteModalOpen(true)}>
                            <Trash2 size={16} />
                            {t('common.delete')}
                        </Button>
                    </div>
                }
            />

            {/* Stats */}
            <HabitStats stats={habitStats} />

            {/* Main Content Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Logs Section (2/3 width) */}
                <div className="lg:col-span-2">
                    <Card className="h-full">
                        <div className="flex items-center justify-between mb-6">
                            <div className="flex items-center gap-3">
                                <div className="p-2.5 bg-primary/5 rounded-lg">
                                    <Activity size={20} className="text-primary" />
                                </div>
                                <div>
                                    <h3 className="text-base font-semibold text-base-content">{t('habits.detail.logs')}</h3>
                                    <p className="text-xs text-base-content/50">
                                        {habitLogs.length} {t('habits.detail.logs').toLowerCase()}
                                    </p>
                                </div>
                            </div>
                            <Button variant="success" size="sm" onClick={() => setIsLogModalOpen(true)}>
                                <CheckCircle2 size={16} />
                                {t('habits.actions.logProgress')}
                            </Button>
                        </div>

                        <HabitLogList logs={habitLogs} onEdit={handleEditLog} onDelete={handleDeleteLog} />

                        <div className="mt-6">
                            <Pagination meta={logsPaging} onPageChange={handlePageChange} />
                        </div>
                    </Card>
                </div>

                {/* Details Sidebar (1/3 width) */}
                <div className="flex flex-col gap-6">
                    {/* Habit Details */}
                    <Card>
                        <div className="flex items-center gap-3 mb-5">
                            <div className="p-2 bg-secondary/5 rounded-lg">
                                <Clock size={18} className="text-secondary" />
                            </div>
                            <h3 className="text-sm font-semibold text-base-content">{t('habits.detail.title')}</h3>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <p className="text-[10px] text-base-content/50 uppercase tracking-wider font-medium mb-1.5">{t('common.active')}</p>
                                <StatusBadge isActive={selectedHabit.is_active} />
                            </div>

                            <div>
                                <p className="text-[10px] text-base-content/50 uppercase tracking-wider font-medium mb-1.5">{t('habits.form.frequency')}</p>
                                <FrequencyBadge frequency={selectedHabit.frequency} />
                            </div>

                            <div>
                                <p className="text-[10px] text-base-content/50 uppercase tracking-wider font-medium mb-1.5">{t('habits.form.targetCount')}</p>
                                <p className="text-base font-semibold text-base-content">
                                    {selectedHabit.target_count}×{' '}
                                    <span className="text-sm text-base-content/50 font-normal">per {t(`habits.frequency.${selectedHabit.frequency}`)}</span>
                                </p>
                            </div>

                            {selectedHabit.description && (
                                <div>
                                    <p className="text-[10px] text-base-content/50 uppercase tracking-wider font-medium mb-1.5">
                                        {t('habits.form.description')}
                                    </p>
                                    <p className="text-sm text-base-content/70 leading-relaxed">{selectedHabit.description}</p>
                                </div>
                            )}
                        </div>
                    </Card>

                    {/* Insights */}
                    <Card>
                        <div className="flex items-center gap-3 mb-5">
                            <div className="p-2 bg-warning/10 rounded-lg">
                                <Lightbulb size={18} className="text-warning" />
                            </div>
                            <h3 className="text-sm font-semibold text-base-content">{t('analytics.completionTrend')}</h3>
                        </div>

                        <div className="space-y-3">
                            <div className="flex items-center gap-3 p-3 bg-base-50 border border-base-200 rounded-lg">
                                <TrendingUp size={18} className="text-success" />
                                <div>
                                    <p className="text-xs font-medium text-base-content/60">{t('habits.stats.currentStreak')}</p>
                                    <p className="text-base font-semibold text-base-content">
                                        {habitStats?.current_streak || 0}{' '}
                                        <span className="text-xs font-normal text-base-content/50">{t('dashboard.days')}</span>
                                    </p>
                                </div>
                            </div>

                            <div className="flex items-center gap-3 p-3 bg-base-50 border border-base-200 rounded-lg">
                                <CheckCircle2 size={18} className="text-primary" />
                                <div>
                                    <p className="text-xs font-medium text-base-content/60">{t('habits.stats.completionRate')}</p>
                                    <p className="text-base font-semibold text-base-content">{habitStats?.completion_rate || 0}%</p>
                                </div>
                            </div>
                        </div>
                    </Card>
                </div>
            </div>

            {/* Modals */}
            <EditHabitModal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)} habit={selectedHabit} />
            <LogHabitModal isOpen={isLogModalOpen} onClose={() => setIsLogModalOpen(false)} habit={selectedHabit} />
            <EditHabitLogModal
                isOpen={isEditLogModalOpen}
                onClose={() => {
                    setIsEditLogModalOpen(false);
                    setSelectedLog(null);
                }}
                log={selectedLog}
            />

            <ConfirmModal
                isOpen={isDeleteModalOpen}
                onClose={() => setIsDeleteModalOpen(false)}
                onConfirm={handleDelete}
                title={t('habits.modal.deleteTitle')}
                message={t('habits.modal.deleteConfirm')}
                confirmText={t('common.delete')}
                variant="danger"
                isLoading={isDeleting}
            />

            <ConfirmModal
                isOpen={isDeleteLogModalOpen}
                onClose={() => setIsDeleteLogModalOpen(false)}
                onConfirm={confirmDeleteLog}
                title={t('habits.detail.deleteLog')}
                message={t('habits.modal.deleteConfirm')}
                confirmText={t('common.delete')}
                variant="danger"
            />
        </div>
    );
}
