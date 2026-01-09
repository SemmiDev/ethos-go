import React, { useEffect, useCallback } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, TouchableOpacity, Dimensions } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useFocusEffect } from '@react-navigation/native';
import { useAuthStore } from '../stores/authStore';
import { useHabitsStore } from '../stores/habitsStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, ProgressBar, Badge } from '../components';
import { Users, CheckCircle2, Flame, TrendingUp, ListTodo, Plus, UserCircle } from 'lucide-react-native';

const { width } = Dimensions.get('window');

const StatCard = ({ title, value, icon, color }) => {
  const { theme } = useThemeStore();

  return (
    <Card style={[styles.statCard, { flex: 1, borderColor: 'transparent', backgroundColor: theme.colors.surface }]}>
      <View style={[styles.iconContainer, { backgroundColor: color + '15' }]}>{icon}</View>
      <View>
        <Text style={[styles.statValue, { color: theme.colors.text }]}>{value}</Text>
        <Text style={[styles.statTitle, { color: theme.colors.textMuted }]}>{title}</Text>
      </View>
    </Card>
  );
};

export default function DashboardScreen({ navigation }) {
  const { user } = useAuthStore();
  const { theme } = useThemeStore();
  const { dashboard, fetchDashboard, isLoading } = useHabitsStore();

  useFocusEffect(
    useCallback(() => {
      fetchDashboard();
    }, [])
  );

  const stats = dashboard?.stats || {
    total_habits: 0,
    completed_today: 0,
    current_streak: 0,
    completion_rate: 0,
  };

  const todayHabits = dashboard?.today_habits || [];

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }} edges={['top']}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        refreshControl={<RefreshControl refreshing={isLoading} onRefresh={fetchDashboard} />}
        showsVerticalScrollIndicator={false}
      >
        {/* Header */}
        <View style={styles.header}>
          <View>
            <Text style={[styles.greeting, { color: theme.colors.textMuted }]}>Good Morning,</Text>
            <Text style={[styles.userName, { color: theme.colors.text }]}>{user?.name}</Text>
          </View>
          <TouchableOpacity onPress={() => navigation.navigate('Settings')}>
            <View style={[styles.avatarPlaceholder, { backgroundColor: theme.colors.primary }]}>
              <Text style={{ color: '#FFF', fontWeight: '600', fontSize: 18 }}>{user?.name?.charAt(0).toUpperCase()}</Text>
            </View>
          </TouchableOpacity>
        </View>

        {/* Stats Grid */}
        <View style={styles.statsGrid}>
          <View style={styles.statsRow}>
            <StatCard title="Habits" value={stats.total_habits} icon={<ListTodo size={24} color={theme.colors.primary} />} color={theme.colors.primary} />
            <View style={{ width: 12 }} />
            <StatCard
              title="Completed"
              value={stats.completed_today}
              icon={<CheckCircle2 size={24} color={theme.colors.success} />}
              color={theme.colors.success}
            />
          </View>
          <View style={{ height: 12 }} />
          <View style={styles.statsRow}>
            <StatCard title="Streak" value={stats.current_streak} icon={<Flame size={24} color={theme.colors.warning} />} color={theme.colors.warning} />
            <View style={{ width: 12 }} />
            <StatCard
              title="Rate"
              value={`${Math.round(stats.completion_rate)}%`}
              icon={<TrendingUp size={24} color={theme.colors.info} />}
              color={theme.colors.info}
            />
          </View>
        </View>

        {/* Today's Habits */}
        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Text style={[styles.sectionTitle, { color: theme.colors.text }]}>Today's Habits</Text>
            <TouchableOpacity onPress={() => navigation.navigate('Habits')} style={{ flexDirection: 'row', alignItems: 'center' }}>
              <Text style={{ color: theme.colors.primary, fontWeight: '600', marginRight: 4 }}>See All</Text>
            </TouchableOpacity>
          </View>

          {todayHabits.length === 0 ? (
            <Card style={styles.emptyCard} padding={true}>
              <Text style={{ textAlign: 'center', color: theme.colors.textMuted, marginBottom: 16 }}>No habits scheduled for today</Text>
              <TouchableOpacity onPress={() => navigation.navigate('Habits')} style={{ flexDirection: 'row', alignItems: 'center', justifyContent: 'center' }}>
                <Plus size={20} color={theme.colors.primary} />
                <Text style={{ color: theme.colors.primary, marginLeft: 8, fontWeight: '600' }}>Add new habit</Text>
              </TouchableOpacity>
            </Card>
          ) : (
            todayHabits.map((habit) => (
              <TouchableOpacity
                key={habit.id}
                onPress={() =>
                  navigation.navigate('Habits', {
                    screen: 'HabitDetail',
                    params: { habitId: habit.id },
                  })
                }
                activeOpacity={0.7}
              >
                <Card style={[styles.habitCard, { backgroundColor: theme.colors.surface }]} padding={false}>
                  <View style={styles.habitCardInner}>
                    <View style={styles.habitHeader}>
                      <View style={{ flex: 1 }}>
                        <Text style={[styles.habitName, { color: theme.colors.text }]}>{habit.name}</Text>
                        <Text style={{ fontSize: 12, color: theme.colors.textMuted, marginTop: 2 }}>
                          {habit.current_count} / {habit.target_count} completed
                        </Text>
                      </View>
                      <Badge text={habit.completed ? 'Done' : 'Pending'} variant={habit.completed ? 'success' : 'neutral'} />
                    </View>
                    <ProgressBar
                      progress={(habit.current_count / habit.target_count) * 100}
                      color={habit.completed ? theme.colors.success : theme.colors.primary}
                      style={{ marginTop: 16 }}
                      height={6}
                    />
                  </View>
                </Card>
              </TouchableOpacity>
            ))
          )}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  scrollContent: {
    padding: 24,
    paddingBottom: 40,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 32,
    marginTop: 8,
  },
  greeting: {
    fontSize: 14,
    marginBottom: 4,
    fontFamily: 'Inter_500Medium',
  },
  userName: {
    fontSize: 24,
    fontFamily: 'Inter_700Bold',
  },
  avatarPlaceholder: {
    width: 44,
    height: 44,
    borderRadius: 22,
    justifyContent: 'center',
    alignItems: 'center',
  },
  statsGrid: {
    marginBottom: 32,
  },
  statsRow: {
    flexDirection: 'row',
  },
  statCard: {
    padding: 16,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  iconContainer: {
    width: 40,
    height: 40,
    borderRadius: 10,
    justifyContent: 'center',
    alignItems: 'center',
  },
  statValue: {
    fontSize: 20,
    fontFamily: 'Inter_700Bold',
    marginBottom: 2,
  },
  statTitle: {
    fontSize: 12,
    fontFamily: 'Inter_500Medium',
  },
  section: {
    marginBottom: 24,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontFamily: 'Inter_600SemiBold',
  },
  habitCard: {
    marginBottom: 12,
    borderWidth: 0,
    shadowOpacity: 0.03,
  },
  habitCardInner: {
    padding: 16,
  },
  habitHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  habitName: {
    fontSize: 16,
    fontFamily: 'Inter_600SemiBold',
  },
  emptyCard: {
    padding: 32,
    alignItems: 'center',
    justifyContent: 'center',
    borderStyle: 'dashed',
    borderWidth: 1,
    backgroundColor: 'transparent',
  },
});
