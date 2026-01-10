import React, { useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, Alert, RefreshControl } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../stores/authStore';
import { useThemeStore } from '../stores/themeStore';
import { Card } from '../components';
import { Trash2, Smartphone, Globe } from 'lucide-react-native';

export default function SessionsScreen() {
  const { theme } = useThemeStore();
  const { sessions, fetchSessions, revokeSession, revokeOtherSessions, isLoading } = useAuthStore();
  console.log('[SessionsScreen] Rendering');

  useEffect(() => {
    console.log('[SessionsScreen] Mounting, fetching sessions...');
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
    if (ua.includes('mobile') || ua.includes('fban') || ua.includes('fbav') || ua.includes('okhttp')) {
      // okhttp is common in android apps
      return <Smartphone size={24} color={theme.colors.textMuted} />;
    }
    return <Globe size={24} color={theme.colors.textMuted} />;
  };

  const renderItem = ({ item }) => {
    try {
      return (
        <View style={{ padding: 16, borderBottomWidth: 1, borderColor: theme.colors.border }}>
          <Text style={{ color: theme.colors.text }}>{item.user_agent || 'Unknown'}</Text>
          <Text style={{ color: theme.colors.textMuted }}>{item.client_ip}</Text>
        </View>
      );
    } catch (e) {
      console.error('Render Item Error', e);
      return null;
    }
  };

  return (
    <View style={[styles.container, { backgroundColor: theme.colors.background }]}>
      <FlatList
        data={sessions}
        renderItem={renderItem}
        keyExtractor={(item) => item.session_id || Math.random().toString()}
        contentContainerStyle={styles.listContent}
        refreshControl={<RefreshControl refreshing={isLoading} onRefresh={fetchSessions} />}
        ListHeaderComponent={
          sessions && sessions.length > 1 ? (
            <View style={{ marginBottom: 16, alignItems: 'flex-end' }}>
              <TouchableOpacity onPress={handleRevokeOthers} style={{ padding: 8 }}>
                <Text style={{ color: theme.colors.error, fontWeight: '600' }}>Sign Out All Other Devices</Text>
              </TouchableOpacity>
            </View>
          ) : null
        }
        ListEmptyComponent={
          !isLoading && (
            <View style={styles.emptyContainer}>
              <Text style={{ color: theme.colors.textMuted }}>No active sessions found.</Text>
            </View>
          )
        }
      />
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
  },
  card: {
    marginBottom: 12,
    padding: 16,
  },
  cardContent: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  iconContainer: {
    width: 40,
    height: 40,
    borderRadius: 20,
    alignItems: 'center',
    justifyContent: 'center',
  },
  infoContainer: {
    flex: 1,
  },
  deviceText: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 2,
  },
  metaText: {
    fontSize: 12,
  },
  badge: {
    alignSelf: 'flex-start',
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 4,
    marginTop: 4,
  },
  badgeText: {
    fontSize: 10,
    fontWeight: '600',
  },
  revokeButton: {
    padding: 8,
  },
  emptyContainer: {
    padding: 40,
    alignItems: 'center',
  },
});
