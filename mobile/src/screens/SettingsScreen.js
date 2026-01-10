import React, { useState } from 'react';
import { useNavigation } from '@react-navigation/native';
import { View, Text, StyleSheet, Switch, TouchableOpacity, ScrollView, Modal, KeyboardAvoidingView, Platform, Alert, Linking } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../stores/authStore';
import { useThemeStore } from '../stores/themeStore';
import { Card, Button, Input } from '../components';
import { User, Moon, Bell, LogOut, ChevronRight, Shield, HelpCircle, Info, X, Key, Trash2, ExternalLink, Mail, MessageCircle } from 'lucide-react-native';

// --- Modals ---

const EditProfileModal = ({ visible, onClose, initialData, onSubmit, isLoading }) => {
  const { theme } = useThemeStore();
  const [name, setName] = useState(initialData?.name || '');
  const [email, setEmail] = useState(initialData?.email || '');

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

const InfoModal = ({ visible, onClose, title, children }) => {
  const { theme } = useThemeStore();

  return (
    <Modal visible={visible} animationType="slide" transparent>
      <View style={styles.modalOverlay}>
        <View style={[styles.modalContent, { backgroundColor: theme.colors.surface, maxHeight: '80%' }]}>
          <View style={styles.modalHeader}>
            <Text style={[styles.modalTitle, { color: theme.colors.text }]}>{title}</Text>
            <TouchableOpacity onPress={onClose}>
              <X size={24} color={theme.colors.textMuted} />
            </TouchableOpacity>
          </View>
          <ScrollView showsVerticalScrollIndicator={false}>{children}</ScrollView>
        </View>
      </View>
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
  const { theme } = useThemeStore();
  const navigation = useNavigation();
  const { isDark, toggleTheme } = useThemeStore();
  const { user, logout, updateProfile, changePassword, deleteAccount, isLoading, sessions, fetchSessions } = useAuthStore();

  const [editProfileVisible, setEditProfileVisible] = useState(false);
  const [changePasswordVisible, setChangePasswordVisible] = useState(false);

  const [helpVisible, setHelpVisible] = useState(false);
  const [privacyVisible, setPrivacyVisible] = useState(false);
  const [aboutVisible, setAboutVisible] = useState(false);

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

  const handleOpenSessions = () => {
    navigation.navigate('Sessions');
  };

  const handleDeleteAccount = () => {
    Alert.alert('Delete Account', 'This action is irreversible. Are you sure you want to delete your account?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: () => {
          Alert.prompt(
            'Confirm Deletion',
            'Please enter your password to confirm.',
            async (password) => {
              if (!password) return;
              const result = await deleteAccount(password);
              if (!result.success) {
                Alert.alert('Error', result.error);
              }
            },
            'secure-text'
          );
        },
      },
    ]);
  };

  const openEmail = () => {
    Linking.openURL('mailto:support@ethos-app.com?subject=Help%20Request');
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
          <Card style={styles.sectionCard} padding={false} noShadow>
            <SettingItem icon={Moon} label="Dark Mode" isSwitch value={isDark} onSwitch={toggleTheme} />
          </Card>
        </View>

        {/* Security */}
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, { color: theme.colors.textMuted }]}>Security</Text>
          <Card style={styles.sectionCard} padding={false} noShadow>
            <SettingItem icon={ExternalLink} label="Active Sessions" onPress={handleOpenSessions} />
            <SettingItem icon={Key} label="Change Password" onPress={() => setChangePasswordVisible(true)} />
          </Card>
        </View>

        {/* Support */}
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, { color: theme.colors.textMuted }]}>Support</Text>
          <Card style={styles.sectionCard} padding={false} noShadow>
            <SettingItem icon={HelpCircle} label="Help Center" onPress={() => setHelpVisible(true)} />
            <SettingItem icon={Shield} label="Privacy Policy" onPress={() => setPrivacyVisible(true)} />
            <SettingItem icon={Info} label="About Ethos" onPress={() => setAboutVisible(true)} />
          </Card>
        </View>

        {/* Account */}
        <View style={styles.section}>
          <Card style={styles.sectionCard} padding={false} noShadow>
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

      {/* Help Center Modal */}
      <InfoModal visible={helpVisible} onClose={() => setHelpVisible(false)} title="Help Center">
        <View style={styles.infoContent}>
          <Text style={[styles.infoHeading, { color: theme.colors.text }]}>Frequently Asked Questions</Text>

          <Text style={[styles.faqQuestion, { color: theme.colors.text }]}>How do I create a new habit?</Text>
          <Text style={[styles.faqAnswer, { color: theme.colors.textMuted }]}>
            Go to the Habits tab and tap the "+" button. Fill in the habit details like name, frequency, and target count.
          </Text>

          <Text style={[styles.faqQuestion, { color: theme.colors.text }]}>How do I track my progress?</Text>
          <Text style={[styles.faqAnswer, { color: theme.colors.textMuted }]}>
            Tap on any habit to view your progress, streaks, and completion history. The Analytics tab shows your overall stats.
          </Text>

          <Text style={[styles.faqQuestion, { color: theme.colors.text }]}>What is vacation mode?</Text>
          <Text style={[styles.faqAnswer, { color: theme.colors.textMuted }]}>
            Vacation mode pauses your streak tracking while you're away. Your streak won't break during this period.
          </Text>

          <Text style={[styles.infoHeading, { color: theme.colors.text, marginTop: 24 }]}>Contact Support</Text>
          <TouchableOpacity style={[styles.contactButton, { backgroundColor: theme.colors.primary }]} onPress={openEmail}>
            <Mail size={18} color="#FFF" />
            <Text style={styles.contactButtonText}>Email Support</Text>
          </TouchableOpacity>
        </View>
      </InfoModal>

      {/* Privacy Policy Modal */}
      <InfoModal visible={privacyVisible} onClose={() => setPrivacyVisible(false)} title="Privacy Policy">
        <View style={styles.infoContent}>
          <Text style={[styles.infoHeading, { color: theme.colors.text }]}>Data Collection</Text>
          <Text style={[styles.infoText, { color: theme.colors.textMuted }]}>
            We collect only the information necessary to provide our habit tracking service. This includes your account information, habit data, and app usage
            statistics.
          </Text>

          <Text style={[styles.infoHeading, { color: theme.colors.text }]}>Data Security</Text>
          <Text style={[styles.infoText, { color: theme.colors.textMuted }]}>
            Your data is encrypted in transit and at rest. We use industry-standard security measures to protect your information.
          </Text>

          <Text style={[styles.infoHeading, { color: theme.colors.text }]}>Your Rights</Text>
          <Text style={[styles.infoText, { color: theme.colors.textMuted }]}>
            You can export or delete your data at any time from the settings. We do not sell your personal information to third parties.
          </Text>

          <Text style={[styles.infoHeading, { color: theme.colors.text }]}>Contact</Text>
          <Text style={[styles.infoText, { color: theme.colors.textMuted }]}>For privacy-related inquiries, please contact privacy@ethos-app.com</Text>
        </View>
      </InfoModal>

      {/* About Modal */}
      <InfoModal visible={aboutVisible} onClose={() => setAboutVisible(false)} title="About Ethos">
        <View style={styles.infoContent}>
          <View style={[styles.aboutLogo, { backgroundColor: theme.colors.primary }]}>
            <Text style={styles.aboutLogoText}>E</Text>
          </View>

          <Text style={[styles.aboutTitle, { color: theme.colors.text }]}>Ethos</Text>
          <Text style={[styles.aboutVersion, { color: theme.colors.textMuted }]}>Version 1.0.0</Text>

          <Text style={[styles.aboutDescription, { color: theme.colors.textMuted }]}>
            Ethos is a habit tracking app designed to help you build positive habits and achieve your goals. Track your daily progress, maintain streaks, and
            visualize your journey to self-improvement.
          </Text>

          <Text style={[styles.infoHeading, { color: theme.colors.text }]}>Features</Text>
          <Text style={[styles.infoText, { color: theme.colors.textMuted }]}>
            • Daily, weekly, and monthly habit tracking{'\n'}• Streak tracking and analytics{'\n'}• Vacation mode to pause tracking{'\n'}• Dark mode support
            {'\n'}• In-app notifications
          </Text>

          <Text style={[styles.aboutCopyright, { color: theme.colors.textMuted }]}>© 2026 Ethos. All rights reserved.</Text>
        </View>
      </InfoModal>
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
  // Info Modal Styles
  infoContent: {
    paddingBottom: 24,
  },
  infoHeading: {
    fontSize: 16,
    fontFamily: 'Inter_600SemiBold',
    marginBottom: 8,
    marginTop: 16,
  },
  infoText: {
    fontSize: 14,
    fontFamily: 'Inter_400Regular',
    lineHeight: 22,
  },
  faqQuestion: {
    fontSize: 15,
    fontFamily: 'Inter_600SemiBold',
    marginTop: 16,
    marginBottom: 4,
  },
  faqAnswer: {
    fontSize: 14,
    fontFamily: 'Inter_400Regular',
    lineHeight: 20,
  },
  contactButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 14,
    borderRadius: 12,
    marginTop: 12,
    gap: 8,
  },
  contactButtonText: {
    color: '#FFF',
    fontSize: 15,
    fontFamily: 'Inter_600SemiBold',
  },
  aboutLogo: {
    width: 80,
    height: 80,
    borderRadius: 20,
    alignItems: 'center',
    justifyContent: 'center',
    alignSelf: 'center',
    marginBottom: 16,
  },
  aboutLogoText: {
    color: '#FFF',
    fontSize: 36,
    fontFamily: 'Inter_700Bold',
  },
  aboutTitle: {
    fontSize: 24,
    fontFamily: 'Inter_700Bold',
    textAlign: 'center',
  },
  aboutVersion: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 16,
  },
  aboutDescription: {
    fontSize: 14,
    fontFamily: 'Inter_400Regular',
    textAlign: 'center',
    lineHeight: 22,
    marginBottom: 24,
  },
  aboutCopyright: {
    fontSize: 12,
    textAlign: 'center',
    marginTop: 24,
  },
  sessionItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
  },
  sessionDevice: {
    fontSize: 16,
    fontFamily: 'Inter_500Medium',
    marginBottom: 4,
  },
  sessionInfo: {
    fontSize: 12,
    fontFamily: 'Inter_400Regular',
  },
  currentBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 8,
  },
});
