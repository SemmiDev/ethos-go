import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, KeyboardAvoidingView, Platform, ScrollView, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../../stores/authStore';
import { useThemeStore } from '../../stores/themeStore';
import { Button, Input } from '../../components';
import { Mail, Lock } from 'lucide-react-native';

export default function LoginScreen({ navigation }) {
  const { theme } = useThemeStore();
  const { login } = useAuthStore();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleLogin = async () => {
    if (!email || !password) {
      setError('Please fill in all fields');
      return;
    }

    setIsLoading(true);
    setError(null);

    const result = await login(email, password);

    setIsLoading(false);

    if (!result.success) {
      setError(result.error);
      Alert.alert('Login Failed', result.error);
    }
  };

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }}>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false}>
          <View style={styles.header}>
            <Text style={[styles.title, { color: theme.colors.text, fontFamily: theme.typography.fontFamily.bold }]}>Welcome back</Text>
            <Text style={[styles.subtitle, { color: theme.colors.textMuted, fontFamily: theme.typography.fontFamily.medium }]}>
              Sign in to continue your habit journey
            </Text>
          </View>

          <View style={styles.form}>
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
              placeholder="Enter your password"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              style={styles.input}
              leftIcon={<Lock size={20} color={theme.colors.textMuted} />}
            />

            <TouchableOpacity onPress={() => navigation.navigate('ForgotPassword')} style={styles.forgotPassword}>
              <Text style={[styles.linkText, { color: theme.colors.primary, fontFamily: theme.typography.fontFamily.medium }]}>Forgot password?</Text>
            </TouchableOpacity>

            {error && <Text style={[styles.errorText, { color: theme.colors.error, fontFamily: theme.typography.fontFamily.medium }]}>{error}</Text>}

            <Button title="Sign In" onPress={handleLogin} loading={isLoading} style={styles.button} />

            <View style={styles.footer}>
              <Text style={[styles.footerText, { color: theme.colors.textMuted, fontFamily: theme.typography.fontFamily.regular }]}>
                Don't have an account?
              </Text>
              <TouchableOpacity onPress={() => navigation.navigate('Register')}>
                <Text style={[styles.linkText, { color: theme.colors.primary, fontFamily: theme.typography.fontFamily.medium }]}>Sign up</Text>
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
  forgotPassword: {
    alignSelf: 'flex-end',
    marginBottom: 24,
  },
  linkText: {
    fontWeight: '600',
  },
  errorText: {
    marginBottom: 16,
    textAlign: 'center',
  },
  button: {
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
