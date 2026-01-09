import React, { useState } from 'react';
import { View, Text, StyleSheet, Switch, TouchableOpacity, ScrollView, Modal, KeyboardAvoidingView, Platform, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../stores/authStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, Button, Input } from '../components';
import { User, Moon, Bell, LogOut, ChevronRight, Shield, HelpCircle, Info, X, Key, Trash2 } from 'lucide-react-native';

// --- Modals ---

const EditProfileModal = ({ visible, onClose, initialData, onSubmit, isLoading }) => {
  const { theme } = useThemeStore();
  const [name, setName] = useState(initialData?.name || '');
  const [email, setEmail] = useState(initialData?.email || '');
  // Timezone could be a picker in a full implementation

  const handleSubmit = () => {
    onSubmit({ name, email });
  };

  return (
    <Modal visible={visible} animationType="slide" transparent>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <View style={styles.modalOverlay}>
          <View style={[styles.modalContent, { backgroundColor: theme.colors.surface }]}>
            <View style={styles.modalHeader}>
              <Text style={[styles.modalTitle, { color: theme.colors.text }]}>Edit Profile</Text>
              <TouchableOpacity onPress={onClose}>
                <X size={24} color={theme.colors.textMuted} />
              </TouchableOpacity>
            </View>
            <Input label="Name" value={name} onChangeText={setName} style={styles.input} />
            <Input label="Email" value={email} onChangeText={setEmail} keyboardType="email-address" style={styles.input} />
            <Button title="Save Changes" onPress={handleSubmit} loading={isLoading} style={styles.modalButton} />
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const ChangePasswordModal = ({ visible, onClose, onSubmit, isLoading }) => {
  const { theme } = useThemeStore();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSubmit = () => {
    if (newPassword !== confirmPassword) {
      Alert.alert('Error', 'New passwords do not match');
      return;
    }
    onSubmit({ current_password: currentPassword, new_password: newPassword });
  };

  return (
    <Modal visible={visible} animationType="slide" transparent>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <View style={styles.modalOverlay}>
          <View style={[styles.modalContent, { backgroundColor: theme.colors.surface }]}>
            <View style={styles.modalHeader}>
              <Text style={[styles.modalTitle, { color: theme.colors.text }]}>Change Password</Text>
              <TouchableOpacity onPress={onClose}>
                <X size={24} color={theme.colors.textMuted} />
              </TouchableOpacity>
            </View>
            <Input label="Current Password" value={currentPassword} onChangeText={setCurrentPassword} secureTextEntry style={styles.input} />
            <Input label="New Password" value={newPassword} onChangeText={setNewPassword} secureTextEntry style={styles.input} />
            <Input label="Confirm New Password" value={confirmPassword} onChangeText={setConfirmPassword} secureTextEntry style={styles.input} />
            <Button title="Update Password" onPress={handleSubmit} loading={isLoading} style={styles.modalButton} />
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

// --- Settings Screen ---

const SettingItem = ({ icon: Icon, label, value, onPress, isSwitch, onSwitch, destructive }) => {
  const { theme } = useThemeStore();
  return (
    <TouchableOpacity
      style={[styles.item, { borderBottomColor: theme.colors.border }]}
      onPress={isSwitch ? onSwitch : onPress}
      disabled={isSwitch}
      activeOpacity={0.7}
    >
      <View style={styles.itemLeft}>
        <View style={[styles.iconContainer, { backgroundColor: destructive ? theme.colors.errorBackground : theme.colors.surfaceVariant }]}>
          <Icon size={20} color={destructive ? theme.colors.error : theme.colors.text} />
        </View>
        <Text style={[styles.itemLabel, { color: destructive ? theme.colors.error : theme.colors.text, fontFamily: theme.typography.fontFamily.medium }]}>
          {label}
        </Text>
      </View>
      {isSwitch ? (
        <Switch
          value={value}
          onValueChange={onSwitch}
          trackColor={{ false: theme.colors.border, true: theme.colors.primary }}
          thumbColor={theme.colors.surface}
          ios_backgroundColor={theme.colors.border}
        />
      ) : (
        <ChevronRight size={20} color={theme.colors.textMuted} />
      )}
    </TouchableOpacity>
  );
};

export default function SettingsScreen() {
  const { theme, isDark, toggleTheme } = useThemeStore();
  const { user, logout, updateProfile, changePassword, deleteAccount, isLoading } = useAuthStore();

  const [editProfileVisible, setEditProfileVisible] = useState(false);
  const [changePasswordVisible, setChangePasswordVisible] = useState(false);

  // Handlers
  const handleUpdateProfile = async (data) => {
    const result = await updateProfile(data);
    if (result.success) {
      setEditProfileVisible(false);
      Alert.alert('Success', 'Profile updated successfully');
    } else {
      Alert.alert('Error', result.error);
    }
  };

  const handleChangePassword = async (data) => {
    const result = await changePassword(data);
    if (result.success) {
      setChangePasswordVisible(false);
      Alert.alert('Success', 'Password updated successfully');
    } else {
      Alert.alert('Error', result.error);
    }
  };

  const handleDeleteAccount = () => {
    Alert.prompt(
      'Delete Account',
      'This action is irreversible. Please enter your password to confirm.',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async (password) => {
            if (!password) return;
            const result = await deleteAccount(password);
            if (!result.success) {
              Alert.alert('Error', result.error);
            }
          },
        },
      ],
      'secure-text'
    );
  };

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }} edges={['top']}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Profile Section */}
        <View style={styles.profileSection}>
          <View style={[styles.avatar, { backgroundColor: theme.colors.primary }]}>
            <Text style={[styles.avatarText, { fontFamily: theme.typography.fontFamily.bold }]}>{user?.name?.charAt(0).toUpperCase() || 'U'}</Text>
          </View>
          <Text style={[styles.profileName, { color: theme.colors.text }]}>{user?.name}</Text>
          <Text style={[styles.profileEmail, { color: theme.colors.textMuted }]}>{user?.email}</Text>
          <View style={{ marginTop: 16 }}>
            <Button
              title="Edit Profile"
              onPress={() => setEditProfileVisible(true)}
              variant="secondary"
              style={{ height: 40, minHeight: 40, paddingVertical: 0 }}
            />
          </View>
        </View>

        {/* Preferences */}
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, { color: theme.colors.textMuted }]}>Preferences</Text>
          <Card style={styles.sectionCard} padding={false}>
            <SettingItem icon={Moon} label="Dark Mode" isSwitch value={isDark} onSwitch={toggleTheme} />
            <SettingItem icon={Bell} label="Notifications" isSwitch value={true} onSwitch={() => {}} />
          </Card>
        </View>

        {/* Security */}
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, { color: theme.colors.textMuted }]}>Security</Text>
          <Card style={styles.sectionCard} padding={false}>
            <SettingItem icon={Key} label="Change Password" onPress={() => setChangePasswordVisible(true)} />
          </Card>
        </View>

        {/* Support */}
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, { color: theme.colors.textMuted }]}>Support</Text>
          <Card style={styles.sectionCard} padding={false}>
            <SettingItem icon={HelpCircle} label="Help Center" onPress={() => {}} />
            <SettingItem icon={Shield} label="Privacy Policy" onPress={() => {}} />
            <SettingItem icon={Info} label="About Ethos" onPress={() => {}} />
          </Card>
        </View>

        {/* Account */}
        <View style={styles.section}>
          <Card style={styles.sectionCard} padding={false}>
            <SettingItem icon={LogOut} label="Log Out" onPress={logout} destructive />
            <SettingItem icon={Trash2} label="Delete Account" onPress={handleDeleteAccount} destructive />
          </Card>
        </View>

        <Text style={[styles.version, { color: theme.colors.textMuted }]}>Version 1.0.0</Text>
      </ScrollView>

      {/* Modals */}
      <EditProfileModal
        visible={editProfileVisible}
        onClose={() => setEditProfileVisible(false)}
        initialData={user}
        onSubmit={handleUpdateProfile}
        isLoading={isLoading}
      />

      <ChangePasswordModal
        visible={changePasswordVisible}
        onClose={() => setChangePasswordVisible(false)}
        onSubmit={handleChangePassword}
        isLoading={isLoading}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  scrollContent: {
    padding: 24,
    paddingBottom: 40,
  },
  profileSection: {
    alignItems: 'center',
    marginBottom: 32,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  avatarText: {
    color: 'white',
    fontSize: 32,
  },
  profileName: {
    fontSize: 20,
    fontFamily: 'Inter_700Bold',
    marginBottom: 4,
  },
  profileEmail: {
    fontSize: 14,
    fontFamily: 'Inter_500Medium',
  },
  section: {
    marginBottom: 24,
  },
  sectionCard: {
    overflow: 'hidden',
    borderWidth: 0,
  },
  sectionTitle: {
    fontSize: 13,
    fontFamily: 'Inter_600SemiBold',
    textTransform: 'uppercase',
    marginBottom: 8,
    marginLeft: 4,
    letterSpacing: 0.5,
  },
  item: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
  },
  itemLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  iconContainer: {
    width: 32,
    height: 32,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  itemLabel: {
    fontSize: 16,
  },
  version: {
    textAlign: 'center',
    fontSize: 12,
    marginTop: 8,
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
  input: {
    marginBottom: 16,
  },
  modalButton: {
    marginTop: 8,
  },
});
