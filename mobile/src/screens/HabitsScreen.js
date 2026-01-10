import React, { useEffect, useState, useMemo } from 'react';
import { View, Text, StyleSheet, FlatList, RefreshControl, TouchableOpacity, Modal, ScrollView, KeyboardAvoidingView, Platform, TextInput } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useHabitsStore } from '../stores/habitsStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, Badge, Button, Input } from '../components';
import { Plus, X, Calendar, Hash, Target, CheckCircle2, Circle, ListTodo, Search, Filter, Pause, Play } from 'lucide-react-native';

const CreateHabitModal = ({ visible, onClose, onSubmit }) => {
  const { theme } = useThemeStore();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [frequency, setFrequency] = useState('daily');
  const [target, setTarget] = useState('1');
  const [reminderTime, setReminderTime] = useState(''); // Format: "HH:MM"

  const handleSubmit = () => {
    if (!name) return;

    // Simple time validation (HH:MM)
    if (reminderTime && !/^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/.test(reminderTime)) {
      alert('Invalid time format. Please use HH:MM (e.g. 08:30)');
      return;
    }

    const payload = {
      name,
      description,
      frequency,
      target_count: parseInt(target) || 1,
    };

    if (reminderTime) {
      payload.reminder_time = reminderTime;
    }

    onSubmit(payload);

    // Reset form
    setName('');
    setDescription('');
    setFrequency('daily');
    setTarget('1');
    setReminderTime('');
    onClose();
  };

  return (
    <Modal visible={visible} animationType="slide" transparent>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <View style={{ flex: 1, backgroundColor: 'rgba(0,0,0,0.5)', justifyContent: 'flex-end' }}>
          <View style={[styles.modalContent, { backgroundColor: theme.colors.surface }]}>
            <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
              <Text style={[styles.modalTitle, { color: theme.colors.text, marginBottom: 0 }]}>New Habit</Text>
              <TouchableOpacity onPress={onClose} style={{ padding: 4 }}>
                <X color={theme.colors.textMuted} size={24} />
              </TouchableOpacity>
            </View>

            <ScrollView showsVerticalScrollIndicator={false}>
              <Input
                label="Habit Name"
                placeholder="e.g., Read 30 mins"
                value={name}
                onChangeText={setName}
                style={styles.modalInput}
                autoCapitalize="sentences"
              />

              <Input
                label="Description (Optional)"
                placeholder="What is this habit about?"
                value={description}
                onChangeText={setDescription}
                style={styles.modalInput}
                autoCapitalize="sentences"
              />

              <View style={{ marginBottom: 24 }}>
                <Text style={[styles.label, { color: theme.colors.text, fontFamily: theme.typography.fontFamily.medium }]}>Frequency</Text>
                <View style={styles.frequencyContainer}>
                  {['daily', 'weekly', 'monthly'].map((f) => (
                    <TouchableOpacity
                      key={f}
                      style={[
                        styles.freqChip,
                        {
                          backgroundColor: frequency === f ? theme.colors.primary : theme.colors.surface,
                          borderColor: frequency === f ? theme.colors.primary : theme.colors.border,
                        },
                      ]}
                      onPress={() => setFrequency(f)}
                    >
                      <Text
                        style={{
                          color: frequency === f ? theme.colors.primaryContent : theme.colors.text,
                          fontWeight: frequency === f ? '600' : '400',
                          textTransform: 'capitalize',
                          fontSize: 14,
                        }}
                      >
                        {f}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </View>
              </View>

              <View style={{ flexDirection: 'row', gap: 16 }}>
                <View style={{ flex: 1 }}>
                  <Input label="Daily Target" placeholder="1" value={target} onChangeText={setTarget} style={styles.modalInput} keyboardType="numeric" />
                </View>
                <View style={{ flex: 1 }}>
                  <Input label="Reminder (HH:MM)" placeholder="08:00" value={reminderTime} onChangeText={setReminderTime} style={styles.modalInput} />
                </View>
              </View>

              <View style={{ marginTop: 8, marginBottom: 24 }}>
                <Button title="Create Habit" onPress={handleSubmit} style={styles.modalButton} />
              </View>
            </ScrollView>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

// Filter Chip Component
const FilterChip = ({ label, active, onPress, theme }) => (
  <TouchableOpacity
    style={[
      styles.filterChip,
      {
        backgroundColor: active ? theme.colors.primary : theme.colors.surface,
        borderColor: active ? theme.colors.primary : theme.colors.border,
      },
    ]}
    onPress={onPress}
  >
    <Text
      style={{
        color: active ? theme.colors.primaryContent : theme.colors.text,
        fontWeight: active ? '600' : '400',
        fontSize: 13,
      }}
    >
      {label}
    </Text>
  </TouchableOpacity>
);

export default function HabitsScreen({ navigation }) {
  const { theme } = useThemeStore();
  const { habits, fetchHabits, createHabit, isLoading, isLoadingMore, pagination } = useHabitsStore();
  const [modalVisible, setModalVisible] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filter, setFilter] = useState('all'); // 'all', 'active', 'inactive'

  // Debounced fetch
  useEffect(() => {
    const timer = setTimeout(() => {
      loadHabits(1);
    }, 300);
    return () => clearTimeout(timer);
  }, [filter, searchQuery]);

  const loadHabits = (page) => {
    const params = {};
    if (filter === 'active') params.active = true;
    if (filter === 'inactive') params.inactive = true;
    if (searchQuery) params.keyword = searchQuery;

    fetchHabits(page, params);
  };

  const handleLoadMore = () => {
    if (pagination.has_next_page && !isLoading && !isLoadingMore) {
      loadHabits(pagination.page + 1);
    }
  };

  const handleCreateHabit = async (data) => {
    await createHabit(data);
  };

  const renderItem = ({ item }) => (
    <TouchableOpacity onPress={() => navigation.navigate('HabitDetail', { habitId: item.id })} activeOpacity={0.7}>
      <Card style={[styles.habitCard, { backgroundColor: theme.colors.surface }]} noShadow>
        <View style={styles.habitContent}>
          <View style={styles.habitInfo}>
            <View style={styles.habitNameRow}>
              <Text style={[styles.habitName, { color: theme.colors.text }]}>{item.name}</Text>
              {!item.is_active && (
                <View style={[styles.statusBadge, { backgroundColor: theme.colors.warningBackground }]}>
                  <Pause size={10} color={theme.colors.warning} />
                  <Text style={[styles.statusText, { color: theme.colors.warning }]}>Paused</Text>
                </View>
              )}
            </View>
            <View style={styles.habitMeta}>
              <View style={styles.metaItem}>
                <Calendar size={12} color={theme.colors.textMuted} />
                <Text style={[styles.habitFreq, { color: theme.colors.textMuted }]}>{item.frequency}</Text>
              </View>
              <View style={styles.metaItem}>
                <Target size={12} color={theme.colors.textMuted} />
                <Text style={[styles.habitFreq, { color: theme.colors.textMuted }]}>Target: {item.target_count}</Text>
              </View>
            </View>
          </View>
          <View style={{ alignItems: 'center', justifyContent: 'center' }}>
            {item.is_active ? (
              <View style={{ backgroundColor: theme.colors.success + '15', borderRadius: 20, padding: 8 }}>
                <CheckCircle2 color={theme.colors.success} size={20} />
              </View>
            ) : (
              <View style={{ backgroundColor: theme.colors.border, borderRadius: 20, padding: 8 }}>
                <Circle color={theme.colors.textMuted} size={20} />
              </View>
            )}
          </View>
        </View>
      </Card>
    </TouchableOpacity>
  );

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }} edges={['top']}>
      {/* Search and Filter */}
      <View style={styles.searchFilterContainer}>
        <View style={[styles.searchContainer, { backgroundColor: theme.colors.surface, borderColor: theme.colors.border }]}>
          <Search size={18} color={theme.colors.textMuted} />
          <TextInput
            style={[styles.searchInput, { color: theme.colors.text }]}
            placeholder="Search habits..."
            placeholderTextColor={theme.colors.textMuted}
            value={searchQuery}
            onChangeText={setSearchQuery}
          />
          {searchQuery.length > 0 && (
            <TouchableOpacity onPress={() => setSearchQuery('')}>
              <X size={18} color={theme.colors.textMuted} />
            </TouchableOpacity>
          )}
        </View>

        <View style={styles.filterContainer}>
          <FilterChip label="All" active={filter === 'all'} onPress={() => setFilter('all')} theme={theme} />
          <FilterChip label="Active" active={filter === 'active'} onPress={() => setFilter('active')} theme={theme} />
          <FilterChip label="Inactive" active={filter === 'inactive'} onPress={() => setFilter('inactive')} theme={theme} />
        </View>
      </View>

      <FlatList
        data={habits}
        renderItem={renderItem}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        refreshControl={<RefreshControl refreshing={isLoading} onRefresh={() => loadHabits(1)} />}
        onEndReached={handleLoadMore}
        onEndReachedThreshold={0.5}
        showsVerticalScrollIndicator={false}
        ListFooterComponent={
          isLoadingMore ? (
            <View style={{ padding: 20, alignItems: 'center' }}>
              <Text style={{ color: theme.colors.textMuted }}>Loading more...</Text>
            </View>
          ) : null
        }
        ListEmptyComponent={
          !isLoading && (
            <View style={styles.emptyContainer}>
              <View style={[styles.emptyIcon, { backgroundColor: theme.colors.surfaceVariant }]}>
                <ListTodo size={48} color={theme.colors.textMuted} />
              </View>
              <Text style={{ color: theme.colors.text, fontSize: 18, fontWeight: '600', marginBottom: 8 }}>
                {searchQuery || filter !== 'all' ? 'No habits found' : 'No habits yet'}
              </Text>
              <Text style={{ color: theme.colors.textMuted, textAlign: 'center', maxWidth: '70%' }}>
                {searchQuery || filter !== 'all' ? 'Try adjusting your search or filter.' : 'Start building your routine by creating your first habit.'}
              </Text>
              {!searchQuery && filter === 'all' && (
                <Button
                  title="Create Habit"
                  onPress={() => setModalVisible(true)}
                  style={{ marginTop: 24, minWidth: 200 }}
                  icon={<Plus color="#FFF" size={20} />}
                />
              )}
            </View>
          )
        }
      />

      {habits.length > 0 && (
        <TouchableOpacity style={[styles.fab, { backgroundColor: theme.colors.primary }]} onPress={() => setModalVisible(true)}>
          <Plus color="#FFF" size={28} />
        </TouchableOpacity>
      )}

      <CreateHabitModal visible={modalVisible} onClose={() => setModalVisible(false)} onSubmit={handleCreateHabit} />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  searchFilterContainer: {
    paddingHorizontal: 24,
    paddingTop: 12,
    paddingBottom: 8,
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: 12,
    paddingHorizontal: 14,
    height: 46,
    gap: 10,
  },
  searchInput: {
    flex: 1,
    fontSize: 15,
    fontFamily: 'Inter_400Regular',
  },
  filterContainer: {
    flexDirection: 'row',
    gap: 8,
    marginTop: 12,
  },
  filterChip: {
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 20,
    borderWidth: 1,
  },
  listContent: {
    padding: 24,
    paddingTop: 8,
    flexGrow: 1,
  },
  habitCard: {
    marginBottom: 10,
    borderWidth: 1,
    padding: 14,
    borderRadius: 12,
  },
  habitContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  habitInfo: {
    flex: 1,
  },
  habitNameRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 6,
  },
  habitName: {
    fontSize: 16,
    fontFamily: 'Inter_600SemiBold',
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
    gap: 3,
  },
  statusText: {
    fontSize: 10,
    fontFamily: 'Inter_500Medium',
  },
  habitMeta: {
    flexDirection: 'row',
    gap: 12,
  },
  metaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  habitFreq: {
    fontSize: 12,
    textTransform: 'capitalize',
    fontFamily: 'Inter_500Medium',
  },
  emptyContainer: {
    padding: 40,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 40,
  },
  emptyIcon: {
    width: 80,
    height: 80,
    borderRadius: 40,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 24,
  },
  modalContent: {
    padding: 24,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    paddingBottom: 48,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: -4,
    },
    shadowOpacity: 0.1,
    shadowRadius: 10,
    elevation: 20,
  },
  modalTitle: {
    fontSize: 20,
    fontFamily: 'Inter_700Bold',
  },
  modalInput: {
    marginBottom: 20,
  },
  label: {
    fontSize: 14,
    marginBottom: 10,
    marginLeft: 4,
  },
  frequencyContainer: {
    flexDirection: 'row',
    gap: 8,
    marginBottom: 0,
  },
  freqChip: {
    flex: 1,
    paddingVertical: 12,
    borderRadius: 12,
    borderWidth: 1,
    alignItems: 'center',
  },
  modalButton: {
    marginBottom: 12,
  },
  fab: {
    position: 'absolute',
    bottom: 24,
    right: 24,
    width: 56,
    height: 56,
    borderRadius: 28,
    alignItems: 'center',
    justifyContent: 'center',
    elevation: 4,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.15,
    shadowRadius: 4,
  },
});
