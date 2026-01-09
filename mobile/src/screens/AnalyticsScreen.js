import React, { useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl, Dimensions, TouchableOpacity } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useHabitsStore } from '../stores/habitsStore';
import { useThemeStore } from '../stores/themeStore';
import { Card } from '../components';
import { Target, TrendingUp, Calendar, Trophy, Badge, CheckCircle2 } from 'lucide-react-native';

const screenWidth = Dimensions.get('window').width;

// Custom Bar Chart Component
const WeeklyBarChart = ({ data, labels, theme }) => {
  const maxVal = Math.max(...data, 100);

  return (
    <View style={styles.chartContainer}>
      <View style={styles.barsContainer}>
        {data.map((value, index) => (
          <View key={index} style={styles.barGroup}>
            <View style={styles.barBackground}>
              <View
                style={[
                  styles.barFill,
                  {
                    height: `${(value / maxVal) * 100}%`,
                    backgroundColor: value > 0 ? theme.colors.primary : theme.colors.primary + '20',
                  },
                ]}
              />
            </View>
            <Text style={[styles.barLabel, { color: theme.colors.textMuted }]}>{labels[index]}</Text>
          </View>
        ))}
      </View>
    </View>
  );
};

const AchievementItem = ({ icon, title, desc, unlocked, theme }) => (
  <View style={[styles.achievementItem, !unlocked && styles.achievementLocked, { borderColor: unlocked ? theme.colors.border : theme.colors.border + '40' }]}>
    <Text style={styles.achievementIcon}>{icon}</Text>
    <Text style={[styles.achievementTitle, { color: unlocked ? theme.colors.text : theme.colors.textMuted }]}>{title}</Text>
    <Text style={[styles.achievementDesc, { color: theme.colors.textMuted }]}>{desc}</Text>
    {unlocked && (
      <View style={[styles.unlockedBadge, { backgroundColor: theme.colors.success + '20' }]}>
        <Text style={[styles.unlockedText, { color: theme.colors.success }]}>Unlocked</Text>
      </View>
    )}
  </View>
);

export default function AnalyticsScreen() {
  const { theme } = useThemeStore();
  const { analytics, fetchAnalytics, isLoading, habits, dashboard } = useHabitsStore();

  useEffect(() => {
    fetchAnalytics();
  }, []);

  const weeklyData = analytics?.weekly_activity?.map((d) => d.completion_rate) || [0, 0, 0, 0, 0, 0, 0];
  const weeklyLabels = analytics?.weekly_activity?.map((d) => d.day.substring(0, 3)) || ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  // Calculate distribution
  const dailyCount = habits.filter((h) => h.frequency === 'daily').length;
  const weeklyCount = habits.filter((h) => h.frequency === 'weekly').length;
  const monthlyCount = habits.filter((h) => h.frequency === 'monthly').length;
  const totalHabits = habits.length || 1;

  const achievements = [
    { icon: '🎯', title: 'First Habit', desc: 'Created your first habit', unlocked: habits.length > 0 },
    { icon: '🔥', title: '7 Day Streak', desc: 'Reached a 7 day streak', unlocked: (dashboard?.longest_streak || 0) >= 7 },
    { icon: '🏆', title: 'Monthly Master', desc: 'Reached a 30 day streak', unlocked: (dashboard?.longest_streak || 0) >= 30 },
    { icon: '⭐', title: '100 Logs', desc: 'Total of 100 activity logs', unlocked: (dashboard?.total_logs || 0) >= 100 },
  ];

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }} edges={['top']}>
      <ScrollView contentContainerStyle={styles.scrollContent} refreshControl={<RefreshControl refreshing={isLoading} onRefresh={fetchAnalytics} />}>
        <Text style={[styles.title, { color: theme.colors.text }]}>Analytics</Text>

        {/* Stats Grid */}
        <View style={styles.grid}>
          <Card style={styles.statsCard}>
            <View style={styles.statsHeader}>
              <Target size={20} color={theme.colors.primary} />
              <Text style={[styles.statsValue, { color: theme.colors.text }]}>{habits.length}</Text>
            </View>
            <Text style={[styles.statsLabel, { color: theme.colors.textMuted }]}>Total Habits</Text>
          </Card>
          <Card style={styles.statsCard}>
            <View style={styles.statsHeader}>
              <TrendingUp size={20} color={theme.colors.success} />
              <Text style={[styles.statsValue, { color: theme.colors.text }]}>{dashboard?.weekly_completion || 0}%</Text>
            </View>
            <Text style={[styles.statsLabel, { color: theme.colors.textMuted }]}>This Week</Text>
          </Card>
          <Card style={styles.statsCard}>
            <View style={styles.statsHeader}>
              <Trophy size={20} color={theme.colors.warning} />
              <Text style={[styles.statsValue, { color: theme.colors.text }]}>{dashboard?.longest_streak || 0}</Text>
            </View>
            <Text style={[styles.statsLabel, { color: theme.colors.textMuted }]}>Best Streak</Text>
          </Card>
          <Card style={styles.statsCard}>
            <View style={styles.statsHeader}>
              <Calendar size={20} color={theme.colors.secondary} />
              <Text style={[styles.statsValue, { color: theme.colors.text }]}>{dashboard?.total_logs || 0}</Text>
            </View>
            <Text style={[styles.statsLabel, { color: theme.colors.textMuted }]}>Total Logs</Text>
          </Card>
        </View>

        {/* Weekly Chart */}
        <Card style={styles.sectionCard}>
          <View style={styles.cardHeader}>
            <View style={[styles.iconBox, { backgroundColor: theme.colors.primary + '15' }]}>
              <TrendingUp size={20} color={theme.colors.primary} />
            </View>
            <View>
              <Text style={[styles.cardTitle, { color: theme.colors.text }]}>Weekly Progress</Text>
              <Text style={[styles.cardSubtitle, { color: theme.colors.textMuted }]}>Completion rate for the last 7 days</Text>
            </View>
          </View>

          <WeeklyBarChart data={weeklyData} labels={weeklyLabels} theme={theme} />
        </Card>

        {/* Habit Distribution */}
        <Card style={styles.sectionCard}>
          <View style={styles.cardHeader}>
            <View style={[styles.iconBox, { backgroundColor: theme.colors.secondary + '15' }]}>
              <Target size={20} color={theme.colors.secondary} />
            </View>
            <Text style={[styles.cardTitle, { color: theme.colors.text }]}>Habit Distribution</Text>
          </View>

          <View style={styles.distributionContainer}>
            <View style={styles.distRow}>
              <View style={styles.distLabelRow}>
                <Text style={[styles.distLabel, { color: theme.colors.text }]}>Daily</Text>
                <Text style={[styles.distCount, { color: theme.colors.textMuted }]}>{dailyCount}</Text>
              </View>
              <View style={[styles.distBarBg, { backgroundColor: theme.colors.border }]}>
                <View style={[styles.distBarFill, { backgroundColor: theme.colors.primary, width: `${(dailyCount / totalHabits) * 100}%` }]} />
              </View>
            </View>

            <View style={styles.distRow}>
              <View style={styles.distLabelRow}>
                <Text style={[styles.distLabel, { color: theme.colors.text }]}>Weekly</Text>
                <Text style={[styles.distCount, { color: theme.colors.textMuted }]}>{weeklyCount}</Text>
              </View>
              <View style={[styles.distBarBg, { backgroundColor: theme.colors.border }]}>
                <View style={[styles.distBarFill, { backgroundColor: theme.colors.secondary, width: `${(weeklyCount / totalHabits) * 100}%` }]} />
              </View>
            </View>

            <View style={styles.distRow}>
              <View style={styles.distLabelRow}>
                <Text style={[styles.distLabel, { color: theme.colors.text }]}>Monthly</Text>
                <Text style={[styles.distCount, { color: theme.colors.textMuted }]}>{monthlyCount}</Text>
              </View>
              <View style={[styles.distBarBg, { backgroundColor: theme.colors.border }]}>
                <View style={[styles.distBarFill, { backgroundColor: theme.colors.warning, width: `${(monthlyCount / totalHabits) * 100}%` }]} />
              </View>
            </View>
          </View>
        </Card>

        {/* Achievements */}
        <Card style={styles.sectionCard}>
          <View style={styles.cardHeader}>
            <View style={[styles.iconBox, { backgroundColor: theme.colors.warning + '15' }]}>
              <Trophy size={20} color={theme.colors.warning} />
            </View>
            <View>
              <Text style={[styles.cardTitle, { color: theme.colors.text }]}>Achievements</Text>
              <Text style={[styles.cardSubtitle, { color: theme.colors.textMuted }]}>milestones unlocked</Text>
            </View>
          </View>

          <View style={styles.achievementsGrid}>
            {achievements.map((item, index) => (
              <AchievementItem key={index} {...item} theme={theme} />
            ))}
          </View>
        </Card>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  scrollContent: {
    padding: 24,
    paddingBottom: 40,
  },
  title: {
    fontSize: 28,
    fontFamily: 'Inter_700Bold',
    marginBottom: 24,
  },
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
    marginBottom: 24,
  },
  statsCard: {
    width: '48%', // Approx 2 columns
    padding: 16,
  },
  statsHeader: {
    flexDirection: 'column',
    alignItems: 'flex-start',
    gap: 8,
    marginBottom: 8,
  },
  statsValue: {
    fontSize: 24,
    fontFamily: 'Inter_700Bold',
  },
  statsLabel: {
    fontSize: 12,
    fontFamily: 'Inter_500Medium',
  },
  sectionCard: {
    padding: 20,
    marginBottom: 24,
  },
  cardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    marginBottom: 20,
  },
  iconBox: {
    width: 40,
    height: 40,
    borderRadius: 12,
    alignItems: 'center',
    justifyContent: 'center',
  },
  cardTitle: {
    fontSize: 18,
    fontFamily: 'Inter_600SemiBold',
  },
  cardSubtitle: {
    fontSize: 12,
  },
  // Chart Styles
  chartContainer: {
    height: 180,
    justifyContent: 'flex-end',
  },
  barsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-end',
    height: 150,
  },
  barGroup: {
    alignItems: 'center',
    gap: 8,
    flex: 1,
  },
  barBackground: {
    width: 8,
    height: '100%',
    justifyContent: 'flex-end',
    borderRadius: 4,
  },
  barFill: {
    width: '100%',
    borderRadius: 4,
  },
  barLabel: {
    fontSize: 12,
    fontFamily: 'Inter_500Medium',
  },
  // Distribution Styles
  distributionContainer: {
    gap: 16,
  },
  distRow: {
    gap: 8,
  },
  distLabelRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  distLabel: {
    fontSize: 14,
    fontFamily: 'Inter_500Medium',
  },
  distCount: {
    fontSize: 14,
  },
  distBarBg: {
    height: 8,
    borderRadius: 4,
    overflow: 'hidden',
  },
  distBarFill: {
    height: '100%',
    borderRadius: 4,
  },
  // Achievement Styles
  achievementsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  achievementItem: {
    width: '48%',
    borderWidth: 1,
    borderRadius: 12,
    padding: 12,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 4,
  },
  achievementLocked: {
    opacity: 0.5,
  },
  achievementIcon: {
    fontSize: 24,
    marginBottom: 8,
  },
  achievementTitle: {
    fontSize: 14,
    fontFamily: 'Inter_600SemiBold',
    textAlign: 'center',
  },
  achievementDesc: {
    fontSize: 10,
    textAlign: 'center',
    marginBottom: 8,
  },
  unlockedBadge: {
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 12,
  },
  unlockedText: {
    fontSize: 10,
    fontFamily: 'Inter_600SemiBold',
  },
});
