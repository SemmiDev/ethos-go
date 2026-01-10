import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, Alert, Modal, KeyboardAvoidingView, Platform, FlatList } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useHabitsStore } from '../stores/habitsStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, Button, ProgressBar, Badge, Input } from '../components';
import { X, Calendar, ClipboardList, CheckCircle2, Edit2, Pause, Play } from 'lucide-react-native';

const LogHabitModal = ({ visible, onClose, habit, onSubmit, isLoading }) => {
  const { theme } = useThemeStore();
  const [count, setCount] = useState('1');
  const [notes, setNotes] = useState('');

  const handleSubmit = () => {
    onSubmit({
      count: parseInt(count) || 1,
      note: notes,
      log_date: new Date().toISOString().split('T')[0],
    });
    setCount('1');
    setNotes('');
  };

  return (
    <Modal visible={visible} animationType="slide" transparent>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <View style={styles.modalOverlay}>
          <View style={[styles.modalContent, { backgroundColor: theme.colors.surface }]}>
            <View style={styles.modalHeader}>
              <Text style={[styles.modalTitle, { color: theme.colors.text }]}>Log Progress</Text>
              <TouchableOpacity onPress={onClose} style={{ padding: 4 }}>
                <X color={theme.colors.textMuted} size={24} />
              </TouchableOpacity>
            </View>

            <Input label="Count" placeholder="1" value={count} onChangeText={setCount} keyboardType="numeric" style={styles.modalInput} />

            <Input
              label="Notes (Optional)"
              placeholder="Any notes about this session?"
              value={notes}
              onChangeText={setNotes}
              style={styles.modalInput}
              multiline
              numberOfLines={3}
            />

            <View style={{ marginTop: 8 }}>
              <Button title="Log Progress" onPress={handleSubmit} loading={isLoading} variant="success" style={styles.modalButton} />
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const EditHabitModal = ({ visible, onClose, habit, onSubmit, isLoading }) => {
  const { theme } = useThemeStore();
  const [name, setName] = useState(habit?.name || '');
  const [description, setDescription] = useState(habit?.description || '');
  const [targetCount, setTargetCount] = useState(habit?.target_count?.toString() || '1');

  // Update state when habit prop changes
  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description || '');
      setTargetCount(habit.target_count?.toString() || '1');
    }
  }, [habit]);

  const handleSubmit = () => {
    onSubmit({
      name,
      description,
      target_count: parseInt(targetCount) || 1,
      // Frequency editing might require backend support or complexity not needed right now
    });
  };

  return (
    <Modal visible={visible} animationType="slide" transparent>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <View style={styles.modalOverlay}>
          <View style={[styles.modalContent, { backgroundColor: theme.colors.surface }]}>
            <View style={styles.modalHeader}>
              <Text style={[styles.modalTitle, { color: theme.colors.text }]}>Edit Habit</Text>
              <TouchableOpacity onPress={onClose} style={{ padding: 4 }}>
                <X color={theme.colors.textMuted} size={24} />
              </TouchableOpacity>
            </View>

            <Input label="Name" value={name} onChangeText={setName} style={styles.modalInput} />
            <Input label="Description" value={description} onChangeText={setDescription} style={styles.modalInput} />
            <Input label="Target Count" value={targetCount} onChangeText={setTargetCount} keyboardType="numeric" style={styles.modalInput} />

            <View style={{ marginTop: 8 }}>
              <Button title="Save Changes" onPress={handleSubmit} loading={isLoading} variant="primary" style={styles.modalButton} />
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

export default function HabitDetailScreen({ route, navigation }) {
  const { habitId } = route.params;
  const { theme } = useThemeStore();
  const { habits, logHabit, updateHabit, deleteHabit, activateHabit, deactivateHabit, fetchHabits, fetchHabitLogs, habitLogs } = useHabitsStore();

  const habit = habits.find((h) => h.id === habitId);
  const [logModalVisible, setLogModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [logLoading, setLogLoading] = useState(false);
  const [editLoading, setEditLoading] = useState(false);
  const [toggleLoading, setToggleLoading] = useState(false);

  useEffect(() => {
    fetchHabitLogs(habitId);
  }, [habitId]);

  if (!habit) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
        <Text>Habit not found</Text>
      </View>
    );
  }

  const handleLogSubmit = async (data) => {
    setLogLoading(true);
    const result = await logHabit(habitId, data);
    setLogLoading(false);

    if (result.success) {
      setLogModalVisible(false);
      fetchHabits();
    } else {
      Alert.alert('Error', result.error);
    }
  };

  const handleEditSubmit = async (data) => {
    setEditLoading(true);
    const result = await updateHabit(habitId, data);
    setEditLoading(false);

    if (result.success) {
      setEditModalVisible(false);
      Alert.alert('Success', 'Habit updated');
    } else {
      Alert.alert('Error', result.error);
    }
  };

  const handleDelete = () => {
    Alert.alert('Delete Habit', 'Are you sure you want to delete this habit?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: async () => {
          const result = await deleteHabit(habitId);
          if (result.success) {
            navigation.goBack();
          } else {
            Alert.alert('Error', result.error);
          }
        },
      },
    ]);
  };

  const handleToggleActive = async () => {
    setToggleLoading(true);
    if (habit.active) {
      const result = await deactivateHabit(habitId);
      if (result.success) {
        Alert.alert('Success', 'Habit paused. Your streak will not be affected.');
      } else {
        Alert.alert('Error', result.error);
      }
    } else {
      const result = await activateHabit(habitId);
      if (result.success) {
        Alert.alert('Success', 'Habit activated!');
      } else {
        Alert.alert('Error', result.error);
      }
    }
    setToggleLoading(false);
  };

  const progress = (habit.current_count / habit.target_count) * 100;

  const renderLogItem = ({ item }) => (
    <View style={[styles.logItem, { borderBottomColor: theme.colors.border }]}>
      <View style={{ flexDirection: 'row', alignItems: 'center', gap: 12 }}>
        <View style={[styles.logIcon, { backgroundColor: theme.colors.success + '15' }]}>
          <CheckCircle2 size={16} color={theme.colors.success} />
        </View>
        <View>
          <Text style={[styles.logCount, { color: theme.colors.text }]}>Completed {item.count}x</Text>
          {item.note && <Text style={[styles.logNote, { color: theme.colors.textMuted }]}>{item.note}</Text>}
        </View>
      </View>
      <Text style={[styles.logDate, { color: theme.colors.textMuted }]}>{new Date(item.log_date || item.completed_at).toLocaleDateString()}</Text>
    </View>
  );

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }} edges={['bottom']}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Header Card */}
        <Card style={styles.headerCard}>
          <View style={styles.headerRow}>
            <View style={{ flex: 1, marginRight: 8 }}>
              <View style={{ flexDirection: 'row', alignItems: 'center', gap: 8 }}>
                <Text style={[styles.title, { color: theme.colors.text }]}>{habit.name}</Text>
                <Badge text={habit.frequency} variant="info" />
              </View>
            </View>
            <TouchableOpacity onPress={() => setEditModalVisible(true)} style={{ padding: 4 }}>
              <Edit2 size={20} color={theme.colors.primary} />
            </TouchableOpacity>
          </View>

          <Text style={[styles.description, { color: theme.colors.textMuted }]}>{habit.description || 'No description provided'}</Text>

          <View style={styles.progressSection}>
            <View style={styles.progressLabels}>
              <Text style={{ color: theme.colors.text }}>Progress Today</Text>
              <Text style={{ color: theme.colors.text, fontWeight: '600' }}>
                {habit.current_count} / {habit.target_count}
              </Text>
            </View>
            <ProgressBar progress={progress} color={progress >= 100 ? theme.colors.success : theme.colors.primary} style={{ height: 12 }} />
          </View>

          <Button
            title={progress >= 100 ? 'Completed for today!' : 'Log Activity'}
            onPress={() => setLogModalVisible(true)}
            disabled={progress >= 100}
            variant={progress >= 100 ? 'secondary' : 'primary'}
            style={{ marginTop: 24 }}
          />
        </Card>

        {/* Stats Row */}
        <View style={styles.statsRow}>
          <Card style={styles.statCard}>
            <Text style={[styles.statValue, { color: theme.colors.text }]}>{habit.streak || 0}</Text>
            <Text style={[styles.statLabel, { color: theme.colors.textMuted }]}>Current Streak</Text>
          </Card>
          <Card style={styles.statCard}>
            <Text style={[styles.statValue, { color: theme.colors.text }]}>{habit.total_completions || 0}</Text>
            <Text style={[styles.statLabel, { color: theme.colors.textMuted }]}>Total Logs</Text>
          </Card>
        </View>

        {/* Activity History */}
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, { color: theme.colors.text }]}>Recent Activity</Text>
          <Card style={styles.historyCard}>
            {habitLogs && habitLogs.length > 0 ? (
              habitLogs.map((log) => <View key={log.id}>{renderLogItem({ item: log })}</View>)
            ) : (
              <Text style={{ color: theme.colors.textMuted, textAlign: 'center', padding: 24 }}>No activity logs yet</Text>
            )}
          </Card>
        </View>

        {/* Actions */}
        <View style={styles.section}>
          <Button
            title={habit.active ? 'Pause Habit' : 'Resume Habit'}
            onPress={handleToggleActive}
            loading={toggleLoading}
            variant="secondary"
            icon={habit.active ? <Pause size={18} color={theme.colors.text} /> : <Play size={18} color={theme.colors.text} />}
            style={{ marginBottom: 12 }}
          />
          <Button title="Delete Habit" onPress={handleDelete} variant="ghost" style={{ borderColor: theme.colors.error, borderWidth: 1 }} />
        </View>
      </ScrollView>

      <LogHabitModal visible={logModalVisible} onClose={() => setLogModalVisible(false)} habit={habit} onSubmit={handleLogSubmit} isLoading={logLoading} />
      <EditHabitModal visible={editModalVisible} onClose={() => setEditModalVisible(false)} habit={habit} onSubmit={handleEditSubmit} isLoading={editLoading} />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  scrollContent: {
    padding: 24,
  },
  headerCard: {
    padding: 24,
    marginBottom: 24,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 12,
  },
  title: {
    fontSize: 24,
    fontFamily: 'Inter_700Bold',
  },
  description: {
    fontSize: 16,
    marginBottom: 24,
    lineHeight: 24,
    fontFamily: 'Inter_400Regular',
  },
  progressSection: {
    gap: 8,
  },
  progressLabels: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 4,
  },
  statsRow: {
    flexDirection: 'row',
    gap: 16,
    marginBottom: 24,
  },
  statCard: {
    flex: 1,
    padding: 16,
    alignItems: 'center',
  },
  statValue: {
    fontSize: 28,
    fontFamily: 'Inter_700Bold',
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 12,
    fontFamily: 'Inter_500Medium',
  },
  section: {
    marginBottom: 24,
  },
  sectionTitle: {
    fontSize: 18,
    fontFamily: 'Inter_600SemiBold',
    marginBottom: 12,
  },
  historyCard: {
    padding: 0,
    overflow: 'hidden',
  },
  logItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
  },
  logIcon: {
    width: 32,
    height: 32,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
  logCount: {
    fontSize: 14,
    fontFamily: 'Inter_600SemiBold',
  },
  logNote: {
    fontSize: 12,
    marginTop: 2,
  },
  logDate: {
    fontSize: 12,
  },
  // Modal Styles
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'flex-end',
  },
  modalContent: {
    padding: 24,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    paddingBottom: 48,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.1,
    shadowRadius: 10,
    elevation: 20,
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 24,
  },
  modalTitle: {
    fontSize: 20,
    fontFamily: 'Inter_700Bold',
  },
  modalInput: {
    marginBottom: 20,
  },
  modalButton: {
    marginBottom: 12,
  },
});
