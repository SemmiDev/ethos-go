import React, { useEffect, useCallback } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, TouchableOpacity, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useFocusEffect } from '@react-navigation/native';
import { useNotificationsStore } from '../stores/notificationsStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, Badge } from '../components';
import { Bell, Check, Trash2, Zap, Calendar, Trophy, Info, CheckCheck } from 'lucide-react-native';
import { formatDistanceToNow } from 'date-fns';

const NotificationIcon = ({ type, theme }) => {
  const iconProps = { size: 20 };

  switch (type) {
    case 'streak_milestone':
      return <Zap {...iconProps} color={theme.colors.warning} />;
    case 'habit_reminder':
      return <Calendar {...iconProps} color={theme.colors.primary} />;
    case 'achievement':
      return <Trophy {...iconProps} color={theme.colors.success} />;
    case 'welcome':
      return <Bell {...iconProps} color={theme.colors.info} />;
    default:
      return <Info {...iconProps} color={theme.colors.info} />;
  }
};

export default function NotificationsScreen({ navigation }) {
  const { theme } = useThemeStore();
  const { notifications, unreadCount, isLoading, error, fetchNotifications, fetchUnreadCount, markAsRead, markAllAsRead, deleteNotification, clearError } =
    useNotificationsStore();

  useFocusEffect(
    useCallback(() => {
      fetchNotifications({ page: 1, per_page: 50 });
      fetchUnreadCount();
    }, [])
  );

  const handleMarkAsRead = async (id) => {
    await markAsRead(id);
  };

  const handleMarkAllAsRead = async () => {
    const result = await markAllAsRead();
    if (!result.success) {
      Alert.alert('Error', result.error);
    }
  };

  const handleDelete = (id) => {
    Alert.alert('Delete Notification', 'Are you sure you want to delete this notification?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: async () => {
          const result = await deleteNotification(id);
          if (!result.success) {
            Alert.alert('Error', result.error);
          }
        },
      },
    ]);
  };

  const handleNotificationPress = (notif) => {
    if (!notif.is_read) {
      markAsRead(notif.id);
    }
    // Navigate to habit detail if notification has habit_id
    if (notif.data?.habit_id) {
      navigation.navigate('Habits', {
        screen: 'HabitDetail',
        params: { habitId: notif.data.habit_id },
      });
    }
  };

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }} edges={['top']}>
      {/* Header */}
      <View style={[styles.header, { borderBottomColor: theme.colors.border }]}>
        <View>
          <Text style={[styles.headerTitle, { color: theme.colors.text }]}>Notifications</Text>
          <Text style={[styles.headerSubtitle, { color: theme.colors.textMuted }]}>{unreadCount > 0 ? `${unreadCount} unread` : 'All caught up!'}</Text>
        </View>
        {unreadCount > 0 && (
          <TouchableOpacity onPress={handleMarkAllAsRead} style={[styles.markAllButton, { backgroundColor: theme.colors.primary + '15' }]}>
            <CheckCheck size={16} color={theme.colors.primary} />
            <Text style={[styles.markAllText, { color: theme.colors.primary }]}>Mark all read</Text>
          </TouchableOpacity>
        )}
      </View>

      <ScrollView
        contentContainerStyle={styles.scrollContent}
        refreshControl={
          <RefreshControl
            refreshing={isLoading}
            onRefresh={() => {
              fetchNotifications({ page: 1, per_page: 50 });
              fetchUnreadCount();
            }}
          />
        }
        showsVerticalScrollIndicator={false}
      >
        {notifications.length === 0 && !isLoading ? (
          <View style={styles.emptyContainer}>
            <View style={[styles.emptyIconContainer, { backgroundColor: theme.colors.surface }]}>
              <Bell size={48} color={theme.colors.textMuted} strokeWidth={1} />
            </View>
            <Text style={[styles.emptyTitle, { color: theme.colors.text }]}>No notifications</Text>
            <Text style={[styles.emptySubtitle, { color: theme.colors.textMuted }]}>You're all caught up! Check back later.</Text>
          </View>
        ) : (
          notifications.map((notif) => (
            <TouchableOpacity key={notif.id} onPress={() => handleNotificationPress(notif)} activeOpacity={0.7}>
              <Card
                style={[
                  styles.notificationCard,
                  {
                    backgroundColor: notif.is_read ? theme.colors.surface : theme.colors.primary + '08',
                    borderColor: notif.is_read ? theme.colors.border : theme.colors.primary + '30',
                  },
                ]}
                padding={false}
              >
                <View style={styles.notificationInner}>
                  <View
                    style={[
                      styles.iconContainer,
                      {
                        backgroundColor: notif.is_read ? theme.colors.background : theme.colors.surface,
                      },
                    ]}
                  >
                    <NotificationIcon type={notif.type} theme={theme} />
                  </View>

                  <View style={styles.contentContainer}>
                    <View style={styles.titleRow}>
                      <Text
                        style={[
                          styles.notificationTitle,
                          {
                            color: notif.is_read ? theme.colors.text : theme.colors.primary,
                            fontFamily: notif.is_read ? 'Inter_500Medium' : 'Inter_600SemiBold',
                          },
                        ]}
                        numberOfLines={1}
                      >
                        {notif.title}
                      </Text>
                      {!notif.is_read && <View style={[styles.unreadDot, { backgroundColor: theme.colors.primary }]} />}
                    </View>

                    <Text style={[styles.notificationMessage, { color: theme.colors.textMuted }]} numberOfLines={2}>
                      {notif.message}
                    </Text>

                    <View style={styles.footerRow}>
                      <Text style={[styles.timestamp, { color: theme.colors.textMuted }]}>
                        {formatDistanceToNow(new Date(notif.created_at), { addSuffix: true })}
                      </Text>

                      <View style={styles.actionButtons}>
                        {!notif.is_read && (
                          <TouchableOpacity
                            onPress={(e) => {
                              e.stopPropagation();
                              handleMarkAsRead(notif.id);
                            }}
                            style={[styles.actionButton, { backgroundColor: theme.colors.primary + '15' }]}
                          >
                            <Check size={14} color={theme.colors.primary} />
                          </TouchableOpacity>
                        )}
                        <TouchableOpacity
                          onPress={(e) => {
                            e.stopPropagation();
                            handleDelete(notif.id);
                          }}
                          style={[styles.actionButton, { backgroundColor: theme.colors.error + '15' }]}
                        >
                          <Trash2 size={14} color={theme.colors.error} />
                        </TouchableOpacity>
                      </View>
                    </View>
                  </View>
                </View>
              </Card>
            </TouchableOpacity>
          ))
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 24,
    paddingVertical: 16,
    borderBottomWidth: 1,
  },
  headerTitle: {
    fontSize: 24,
    fontFamily: 'Inter_700Bold',
  },
  headerSubtitle: {
    fontSize: 14,
    fontFamily: 'Inter_400Regular',
    marginTop: 2,
  },
  markAllButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 8,
    gap: 6,
  },
  markAllText: {
    fontSize: 13,
    fontFamily: 'Inter_600SemiBold',
  },
  scrollContent: {
    padding: 16,
    paddingBottom: 40,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 80,
  },
  emptyIconContainer: {
    width: 100,
    height: 100,
    borderRadius: 50,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 20,
  },
  emptyTitle: {
    fontSize: 18,
    fontFamily: 'Inter_600SemiBold',
    marginBottom: 8,
  },
  emptySubtitle: {
    fontSize: 14,
    fontFamily: 'Inter_400Regular',
    textAlign: 'center',
    paddingHorizontal: 40,
  },
  notificationCard: {
    marginBottom: 10,
  },
  notificationInner: {
    flexDirection: 'row',
    padding: 16,
  },
  iconContainer: {
    width: 44,
    height: 44,
    borderRadius: 12,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 12,
  },
  contentContainer: {
    flex: 1,
  },
  titleRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 4,
  },
  notificationTitle: {
    fontSize: 15,
    flex: 1,
  },
  unreadDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginLeft: 8,
  },
  notificationMessage: {
    fontSize: 13,
    fontFamily: 'Inter_400Regular',
    lineHeight: 18,
    marginBottom: 8,
  },
  footerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  timestamp: {
    fontSize: 12,
    fontFamily: 'Inter_400Regular',
  },
  actionButtons: {
    flexDirection: 'row',
    gap: 8,
  },
  actionButton: {
    width: 28,
    height: 28,
    borderRadius: 14,
    alignItems: 'center',
    justifyContent: 'center',
  },
});
