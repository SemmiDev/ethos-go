import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Check, Trash2, Bell, Calendar, Trophy, Zap, Info } from 'lucide-react';
import { useNotificationStore } from '../../stores/notificationStore';
import { formatDistanceToNow } from 'date-fns';
import { Header } from '../../components/layout/Sidebar';

export function NotificationsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { notifications, unreadCount, fetchNotifications, markAsRead, markAllAsRead, deleteNotification, isLoading } = useNotificationStore();

  useEffect(() => {
    fetchNotifications({ page: 1, per_page: 50 });
  }, [fetchNotifications]);

  const getIcon = (type) => {
    switch (type) {
      case 'streak_milestone':
        return <Zap className="text-warning" size={20} />;
      case 'habit_reminder':
        return <Calendar className="text-primary" size={20} />;
      case 'achievement':
        return <Trophy className="text-success" size={20} />;
      default:
        return <Info className="text-info" size={20} />;
    }
  };

  return (
    <div className="p-4 md:p-6 max-w-4xl mx-auto">
      <Header
        title={t('notifications.title')}
        subtitle={t('notifications.subtitle')}
        actions={
          unreadCount > 0 && (
            <button onClick={() => markAllAsRead()} className="btn btn-sm btn-ghost text-primary">
              {t('notifications.markAllRead')}
            </button>
          )
        }
      />

      <div className="space-y-3">
        {isLoading && notifications.length === 0 ? (
          <div className="flex justify-center py-12">
            <span className="loading loading-spinner loading-lg text-primary"></span>
          </div>
        ) : notifications.length === 0 ? (
          <div className="text-center py-20 bg-base-100 rounded-xl border border-base-200">
            <div className="w-16 h-16 bg-base-200 rounded-full flex items-center justify-center mx-auto mb-4">
              <Bell size={32} className="opacity-20" />
            </div>
            <h3 className="font-semibold text-lg">{t('notifications.noNotifications')}</h3>
            <p className="text-base-content/50">{t('notifications.noNotificationsDesc')}</p>
          </div>
        ) : (
          notifications.map((notif) => (
            <div
              key={notif.id}
              className={`
                    relative group p-4 md:p-5 rounded-xl border transition-all duration-200 cursor-pointer
                    ${notif.is_read ? 'bg-base-100 border-base-200' : 'bg-primary/5 border-primary/20'}
                    hover:border-primary/30 hover:shadow-sm
                `}
              onClick={() => {
                if (!notif.is_read) markAsRead(notif.id);
                if (notif.data?.habit_id) {
                  navigate(`/habits/${notif.data.habit_id}`);
                }
              }}
            >
              <div className="flex gap-4">
                <div className={`mt-1 p-3 rounded-xl h-fit ${notif.is_read ? 'bg-base-200' : 'bg-base-100 shadow-sm'}`}>{getIcon(notif.type)}</div>
                <div className="flex-1 min-w-0">
                  <div className="flex justify-between items-start mb-1 gap-2">
                    <h3 className={`text-base font-semibold ${notif.is_read ? 'text-base-content' : 'text-primary'}`}>{notif.title}</h3>
                    <div className="flex items-center gap-3 shrink-0">
                      <span className="text-xs text-base-content/40 whitespace-nowrap">
                        {formatDistanceToNow(new Date(notif.created_at), { addSuffix: true })}
                      </span>
                      {/* Actions */}
                      <div className="flex gap-1" onClick={(e) => e.stopPropagation()}>
                        {!notif.is_read && (
                          <button
                            onClick={() => markAsRead(notif.id)}
                            className="btn btn-circle btn-xs btn-ghost text-primary hover:bg-primary/10"
                            title={t('notifications.markAllRead')}
                          >
                            <Check size={14} />
                          </button>
                        )}
                        <button
                          onClick={() => deleteNotification(notif.id)}
                          className="btn btn-circle btn-xs btn-ghost text-error hover:bg-error/10"
                          title={t('common.delete')}
                        >
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </div>
                  </div>
                  <p className="text-sm md:text-base text-base-content/70 leading-relaxed group-hover:text-base-content transition-colors">{notif.message}</p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
