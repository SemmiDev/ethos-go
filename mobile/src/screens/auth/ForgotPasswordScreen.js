import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, KeyboardAvoidingView, Platform, ScrollView, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useAuthStore } from '../../stores/authStore';
import { useThemeStore } from '../../stores/themeStore';
import { Button, Input } from '../../components';
import { authAPI } from '../../api/auth';

export default function ForgotPasswordScreen({ navigation }) {
  const { theme } = useThemeStore();

  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  const handleSubmit = async () => {
    if (!email) {
      Alert.alert('Error', 'Please enter your email address');
      return;
    }

    setIsLoading(true);
    try {
      await authAPI.forgotPassword({ email });
      setIsSuccess(true);
    } catch (error) {
      Alert.alert('Error', error.response?.data?.message || 'Failed to send reset email');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: theme.colors.background }}>
      <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'} style={{ flex: 1 }}>
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
            <Text style={{ color: theme.colors.textMuted }}>← Back to Login</Text>
          </TouchableOpacity>

          <View style={styles.header}>
            <Text style={[styles.title, { color: theme.colors.text }]}>{isSuccess ? 'Check your email' : 'Forgot Password'}</Text>
            <Text style={[styles.subtitle, { color: theme.colors.textMuted }]}>
              {isSuccess ? `We've sent password reset instructions to ${email}` : "Enter your email address and we'll send you a link to reset your password."}
            </Text>
          </View>

          {!isSuccess && (
            <View style={styles.form}>
              <Input
                placeholder="Email address"
                value={email}
                onChangeText={setEmail}
                keyboardType="email-address"
                autoCapitalize="none"
                style={styles.input}
              />

              <Button title="Send Reset Link" onPress={handleSubmit} loading={isLoading} style={styles.button} />
            </View>
          )}

          {isSuccess && <Button title="Back to Login" onPress={() => navigation.navigate('Login')} variant="secondary" style={styles.button} />}
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
  backButton: {
    position: 'absolute',
    top: 24,
    left: 24,
    zIndex: 1,
  },
  header: {
    marginBottom: 40,
    marginTop: 60,
  },
  title: {
    fontSize: 32,
    fontWeight: '700',
    marginBottom: 16,
  },
  subtitle: {
    fontSize: 16,
    lineHeight: 24,
  },
  form: {
    width: '100%',
  },
  input: {
    marginBottom: 24,
  },
  button: {
    marginBottom: 24,
  },
});
