import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Bell } from 'lucide-react';
import { useNotificationStore } from '../../stores/notificationStore';

export function NotificationBell() {
  const navigate = useNavigate();
  const { unreadCount, fetchUnreadCount } = useNotificationStore();

  useEffect(() => {
    // Initial fetch
    fetchUnreadCount();

    // Poll every minute
    const interval = setInterval(fetchUnreadCount, 60000);
    return () => clearInterval(interval);
  }, [fetchUnreadCount]);

  return (
    <button onClick={() => navigate('/notifications')} className="btn btn-ghost btn-circle relative">
      <Bell size={20} />
      {unreadCount > 0 && <span className="absolute top-2 right-2 w-2.5 h-2.5 bg-error rounded-full ring-2 ring-base-100 animate-pulse" />}
    </button>
  );
}
