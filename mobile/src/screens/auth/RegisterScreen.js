import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, KeyboardAvoidingView, Platform, ScrollView, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../../stores/authStore';
import { useThemeStore } from '../../stores/themeStore';
import { Button, Input } from '../../components';
import { Mail, Lock, User } from 'lucide-react-native';

export default function RegisterScreen({ navigation }) {
  const { theme } = useThemeStore();
  const { register, login } = useAuthStore();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleRegister = async () => {
    if (!name || !email || !password) {
      setError('Please fill in all fields');
      return;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setIsLoading(true);
    setError(null);

    // Get device timezone
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';

    const result = await register(name, email, password, timezone);

    if (result.success) {
      // Auto login after registration
      const loginResult = await login(email, password);
      setIsLoading(false);

      if (!loginResult.success) {
        Alert.alert('Registration Successful', 'Please sign in with your new account.');
        navigation.navigate('Login');
      }
    } else {
      setIsLoading(false);
      setError(result.error);
      Alert.alert('Registration Failed', result.error);
    }
  };

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }}>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false}>
          <View style={styles.header}>
            <Text style={[styles.title, { color: theme.colors.text, fontFamily: theme.typography.fontFamily.bold }]}>Create account</Text>
            <Text style={[styles.subtitle, { color: theme.colors.textMuted, fontFamily: theme.typography.fontFamily.medium }]}>
              Start tracking your habits today
            </Text>
          </View>

          <View style={styles.form}>
            <Input
              label="Full Name"
              placeholder="Enter your name"
              value={name}
              onChangeText={setName}
              autoCapitalize="words"
              style={styles.input}
              leftIcon={<User size={20} color={theme.colors.textMuted} />}
            />

            <Input
              label="Email address"
              placeholder="Enter your email"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              style={styles.input}
              leftIcon={<Mail size={20} color={theme.colors.textMuted} />}
            />

            <Input
              label="Password"
              placeholder="Create a password"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              style={styles.input}
              leftIcon={<Lock size={20} color={theme.colors.textMuted} />}
            />

            <Input
              label="Confirm Password"
              placeholder="Repeat your password"
              value={confirmPassword}
              onChangeText={setConfirmPassword}
              secureTextEntry
              style={styles.input}
              leftIcon={<Lock size={20} color={theme.colors.textMuted} />}
            />

            {error && <Text style={[styles.errorText, { color: theme.colors.error, fontFamily: theme.typography.fontFamily.medium }]}>{error}</Text>}

            <Button title="Create Account" onPress={handleRegister} loading={isLoading} style={styles.button} />

            <View style={styles.footer}>
              <Text style={[styles.footerText, { color: theme.colors.textMuted, fontFamily: theme.typography.fontFamily.regular }]}>
                Already have an account?
              </Text>
              <TouchableOpacity onPress={() => navigation.navigate('Login')}>
                <Text style={[styles.linkText, { color: theme.colors.primary, fontFamily: theme.typography.fontFamily.medium }]}>Sign in</Text>
              </TouchableOpacity>
            </View>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  scrollContent: {
    flexGrow: 1,
    padding: 24,
    justifyContent: 'center',
  },
  header: {
    marginBottom: 40,
  },
  title: {
    fontSize: 28,
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    lineHeight: 24,
  },
  form: {
    width: '100%',
  },
  input: {
    marginBottom: 20,
  },
  linkText: {
    fontWeight: '600',
  },
  errorText: {
    marginBottom: 16,
    textAlign: 'center',
  },
  button: {
    marginTop: 8,
    marginBottom: 24,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: 8,
  },
  footerText: {
    fontSize: 14,
  },
});
