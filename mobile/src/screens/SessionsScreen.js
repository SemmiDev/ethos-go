import React, { useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, Alert, RefreshControl } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../stores/authStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, Button } from '../components';
import { Trash2, Smartphone, Globe, Monitor, Laptop, LogOut } from 'lucide-react-native';

export default function SessionsScreen() {
  const { theme } = useThemeStore();
  const { sessions, fetchSessions, revokeSession, revokeOtherSessions, isSessionsLoading } = useAuthStore();

  useEffect(() => {
    fetchSessions();
  }, []);

  const handleRevokeOthers = () => {
    Alert.alert('Sign Out All Other Devices', 'Are you sure you want to sign out from all other devices?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Sign Out',
        style: 'destructive',
        onPress: async () => {
          const result = await revokeOtherSessions();
          if (result.success) {
            Alert.alert('Success', 'Signed out from other devices.');
          } else {
            Alert.alert('Error', result.error);
          }
        },
      },
    ]);
  };

  const handleRevoke = (sessionId) => {
    Alert.alert('Revoke Session', 'Are you sure you want to revoke this session?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Revoke',
        style: 'destructive',
        onPress: async () => {
          const result = await revokeSession(sessionId);
          if (!result.success) {
            Alert.alert('Error', result.error);
          }
        },
      },
    ]);
  };

  const getDeviceIcon = (userAgent) => {
    if (!userAgent) return <Globe size={24} color={theme.colors.textMuted} />;
    const ua = userAgent.toLowerCase();
    if (ua.includes('mobile') || ua.includes('android') || ua.includes('iphone')) {
      return <Smartphone size={24} color={theme.colors.textMuted} />;
    }
    if (ua.includes('macintosh') || ua.includes('windows') || ua.includes('linux')) {
      return <Laptop size={24} color={theme.colors.textMuted} />;
    }
    return <Globe size={24} color={theme.colors.textMuted} />;
  };

  const renderItem = ({ item }) => {
    const isCurrent = item.is_current;

    return (
      <View style={[styles.sessionItem, { borderColor: theme.colors.border, backgroundColor: theme.colors.surface }]}>
        <View style={styles.itemLeft}>
          <View style={[styles.iconContainer, { backgroundColor: theme.colors.surfaceVariant }]}>{getDeviceIcon(item.user_agent)}</View>
          <View style={styles.sessionInfo}>
            <Text style={[styles.deviceName, { color: theme.colors.text }]} numberOfLines={1}>
              {item.user_agent || 'Unknown Device'}
            </Text>
            <Text style={[styles.deviceMeta, { color: theme.colors.textMuted }]}>{item.client_ip}</Text>
          </View>
        </View>

        <View style={styles.itemRight}>
          <View style={[styles.badge, { backgroundColor: isCurrent ? theme.colors.primary + '20' : '#22c55e20' }]}>
            <Text style={[styles.badgeText, { color: isCurrent ? theme.colors.primary : '#22c55e' }]}>{isCurrent ? 'Current' : 'Active'}</Text>
          </View>

          {!isCurrent && (
            <TouchableOpacity onPress={() => handleRevoke(item.session_id)} style={[styles.revokeButton, { backgroundColor: theme.colors.error + '20' }]}>
              <Trash2 size={16} color={theme.colors.error} />
            </TouchableOpacity>
          )}
        </View>
      </View>
    );
  };

  const hasOtherSessions = sessions && sessions.length > 1;

  return (
    <View style={[styles.container, { backgroundColor: theme.colors.background }]}>
      <FlatList
        data={sessions}
        renderItem={renderItem}
        keyExtractor={(item) => item.session_id || Math.random().toString()}
        contentContainerStyle={styles.listContent}
        refreshControl={<RefreshControl refreshing={isSessionsLoading} onRefresh={fetchSessions} />}
        ListEmptyComponent={
          !isSessionsLoading && (
            <View style={styles.emptyContainer}>
              <Text style={{ color: theme.colors.textMuted }}>No active sessions found.</Text>
            </View>
          )
        }
      />

      {hasOtherSessions && (
        <View style={[styles.footer, { backgroundColor: theme.colors.surface, borderTopColor: theme.colors.border }]}>
          <Button
            title="Sign Out All Other Devices"
            onPress={handleRevokeOthers}
            variant="destructive"
            style={styles.signOutButton}
            icon={<LogOut size={18} color="#FFF" />}
          />
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  listContent: {
    padding: 16,
    paddingBottom: 32,
    gap: 12,
  },
  sessionItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderRadius: 12,
    borderWidth: 1,
    gap: 12,
  },
  itemLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    flex: 1,
  },
  itemRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  iconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    alignItems: 'center',
    justifyContent: 'center',
  },
  sessionInfo: {
    flex: 1,
  },
  deviceName: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 2,
  },
  deviceMeta: {
    fontSize: 12,
  },
  badge: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
  },
  badgeText: {
    fontSize: 12,
    fontWeight: '600',
  },
  revokeButton: {
    width: 32,
    height: 32,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
  emptyContainer: {
    padding: 40,
    alignItems: 'center',
  },
  footer: {
    padding: 16,
    borderTopWidth: 1,
    paddingBottom: 32, // Extra padding for safe area
  },
  signOutButton: {
    width: '100%',
  },
});
