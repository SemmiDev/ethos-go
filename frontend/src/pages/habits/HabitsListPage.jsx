import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Plus, Search, Grid3X3, List, Target } from 'lucide-react';
import { Header } from '../../components/layout/Sidebar';
import { Button } from '../../components/ui/Button';
import { Pagination } from '../../components/ui/Pagination';
import { PageLoader } from '../../components/ui/Loading';
import { ConfirmModal } from '../../components/ui/Modal';
import { HabitCard, CreateHabitModal, EditHabitModal, LogHabitModal } from '../../components/habits';
import { useHabitsStore } from '../../stores/habitsStore';
import { useUIStore } from '../../stores/uiStore';

export function HabitsListPage() {
  const { t } = useTranslation();
  const { habits, isLoading, paging, filter, searchQuery, fetchHabits, deleteHabit, setPage, setFilter, setSearchQuery } = useHabitsStore();
  const { addToast } = useUIStore();

  const [viewMode, setViewMode] = useState('grid');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState(null);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    fetchHabits({ page: 1, per_page: 9 });
  }, [fetchHabits]);

  const handleEdit = (habit) => {
    setSelectedHabit(habit);
    setIsEditModalOpen(true);
  };

  const handleLog = (habit) => {
    setSelectedHabit(habit);
    setIsLogModalOpen(true);
  };

  const handleDeleteClick = (habit) => {
    setSelectedHabit(habit);
    setIsDeleteModalOpen(true);
  };

  const handleDeleteConfirm = async () => {
    setIsDeleting(true);
    const result = await deleteHabit(selectedHabit.id);
    setIsDeleting(false);

    if (result.success) {
      addToast({ type: 'success', title: t('common.success'), message: t('toast.habitDeleted') });
      setIsDeleteModalOpen(false);
    } else {
      addToast({ type: 'error', title: t('common.error'), message: result.error });
    }
  };

  if (isLoading && !habits.length) {
    return <PageLoader />;
  }

  return (
    <div className="space-y-6">
      <Header
        title={t('habits.title')}
        subtitle={t('habits.subtitle')}
        actions={
          <Button variant="primary" onClick={() => setIsCreateModalOpen(true)}>
            <Plus size={16} />
            {t('habits.createHabit')}
          </Button>
        }
      />

      {/* Filters & Controls */}
      <div className="bg-base-100 border border-base-300 rounded-lg p-4">
        <div className="flex flex-col sm:flex-row gap-4 justify-between items-stretch sm:items-center">
          {/* Search */}
          <div className="flex-1 max-w-md">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
              <input
                type="text"
                className="
                                    w-full pl-10 pr-4 py-2.5 text-sm
                                    bg-base-100 border border-base-300 rounded-md
                                    placeholder:text-base-content/40
                                    focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20
                                "
                placeholder={t('habits.searchPlaceholder')}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          <div className="flex items-center gap-3">
            {/* Filter */}
            <div className="relative">
              <select
                className="
                                    pl-3 pr-8 py-2.5 text-sm
                                    bg-base-100 border border-base-300 rounded-md
                                    appearance-none cursor-pointer
                                    focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20
                                "
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
              >
                <option value="all">{t('habits.filters.all')}</option>
                <option value="active">{t('habits.filters.active')}</option>
                <option value="inactive">{t('habits.filters.inactive')}</option>
              </select>
              <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
                <svg className="h-4 w-4 text-base-content/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>

            {/* View Toggle */}
            <div className="hidden sm:flex border border-base-300 rounded-md overflow-hidden">
              <button
                onClick={() => setViewMode('grid')}
                className={`p-2.5 ${viewMode === 'grid' ? 'bg-primary text-primary-content' : 'bg-base-100 text-base-content/60 hover:bg-base-200'}`}
              >
                <Grid3X3 size={16} />
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`p-2.5 border-l border-base-300 ${
                  viewMode === 'list' ? 'bg-primary text-primary-content' : 'bg-base-100 text-base-content/60 hover:bg-base-200'
                }`}
              >
                <List size={16} />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Habits Grid/List */}
      {habits.length === 0 ? (
        <div className="bg-base-100 border border-base-300 rounded-lg text-center py-16">
          <div className="flex flex-col items-center justify-center">
            <div className="p-5 bg-base-200 rounded-full mb-4">
              <Target size={32} className="text-base-content/30" />
            </div>
            <h3 className="text-base font-semibold text-base-content mb-1">{t('habits.noHabits')}</h3>
            <p className="text-sm text-base-content/50 mb-6 max-w-sm mx-auto">
              {searchQuery || filter !== 'all' ? t('habits.noHabitsDesc') : t('habits.noHabitsDesc')}
            </p>
            {!searchQuery && filter === 'all' && (
              <Button variant="primary" onClick={() => setIsCreateModalOpen(true)}>
                <Plus size={16} />
                {t('dashboard.createFirst')}
              </Button>
            )}
          </div>
        </div>
      ) : (
        <>
          <div className={viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5' : 'flex flex-col gap-4'}>
            {habits.map((habit) => (
              <HabitCard
                key={habit.id}
                habit={habit}
                onLog={() => handleLog(habit)}
                onEdit={() => handleEdit(habit)}
                onDelete={() => handleDeleteClick(habit)}
              />
            ))}
          </div>

          <Pagination meta={paging} onPageChange={setPage} />
        </>
      )}

      {/* Modals */}
      <CreateHabitModal isOpen={isCreateModalOpen} onClose={() => setIsCreateModalOpen(false)} />
      <EditHabitModal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)} habit={selectedHabit} />
      <LogHabitModal isOpen={isLogModalOpen} onClose={() => setIsLogModalOpen(false)} habit={selectedHabit} />
      <ConfirmModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title={t('habits.modal.deleteTitle')}
        message={t('habits.modal.deleteConfirm')}
        confirmText={t('common.delete')}
        variant="danger"
        isLoading={isDeleting}
      />
    </div>
  );
}
