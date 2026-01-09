import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator, TextInput as RNTextInput, Animated } from 'react-native';
import { useThemeStore } from '../stores/themeStore';

// Primary Button
export const Button = ({ title, onPress, variant = 'primary', disabled = false, loading = false, style, icon }) => {
  const { theme } = useThemeStore();

  const getButtonStyle = () => {
    switch (variant) {
      case 'secondary':
        return {
          backgroundColor: theme.colors.surface,
          borderWidth: 1,
          borderColor: theme.colors.border,
        };
      case 'ghost':
        return {
          backgroundColor: 'transparent',
          borderWidth: 0,
        };
      case 'danger':
        return {
          backgroundColor: theme.colors.errorBackground,
          borderWidth: 0,
        };
      default:
        return {
          backgroundColor: theme.colors.primary,
        };
    }
  };

  const getTextColor = () => {
    switch (variant) {
      case 'secondary':
      case 'ghost':
        return theme.colors.text;
      case 'danger':
        return theme.colors.error;
      default:
        return theme.colors.primaryContent;
    }
  };

  return (
    <TouchableOpacity
      onPress={onPress}
      disabled={disabled || loading}
      style={[styles.button, getButtonStyle(), disabled && styles.disabled, style]}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator color={getTextColor()} size="small" />
      ) : (
        <View style={styles.buttonContent}>
          {icon && <View style={{ marginRight: 8 }}>{icon}</View>}
          <Text style={[styles.buttonText, { color: getTextColor(), fontFamily: theme.typography.fontFamily.medium }]}>{title}</Text>
        </View>
      )}
    </TouchableOpacity>
  );
};

// Card Component
export const Card = ({ children, style, padding = true }) => {
  const { theme } = useThemeStore();

  return (
    <View
      style={[
        styles.card,
        {
          backgroundColor: theme.colors.surface,
          borderColor: theme.colors.border,
          shadowColor: theme.colors.text, // Use text color for shadow to adapt to dark mode
          padding: padding ? 16 : 0,
        },
        style,
      ]}
    >
      {children}
    </View>
  );
};

// Input Component
export const Input = ({
  value,
  onChangeText,
  placeholder,
  secureTextEntry = false,
  keyboardType = 'default',
  autoCapitalize = 'none',
  error,
  style,
  leftIcon,
  rightIcon,
  label,
}) => {
  const { theme } = useThemeStore();
  const [isFocused, setIsFocused] = useState(false);

  return (
    <View style={[styles.inputContainer, style]}>
      {label && <Text style={[styles.label, { color: theme.colors.text, fontFamily: theme.typography.fontFamily.medium }]}>{label}</Text>}
      <View
        style={[
          styles.inputWrapper,
          {
            backgroundColor: theme.colors.surface,
            borderColor: error ? theme.colors.error : isFocused ? theme.colors.primary : theme.colors.border,
          },
        ]}
      >
        {leftIcon && <View style={styles.leftIcon}>{leftIcon}</View>}
        <RNTextInput
          value={value}
          onChangeText={onChangeText}
          placeholder={placeholder}
          placeholderTextColor={theme.colors.textMuted}
          secureTextEntry={secureTextEntry}
          keyboardType={keyboardType}
          autoCapitalize={autoCapitalize}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          style={[
            styles.textInput,
            {
              color: theme.colors.text,
              fontFamily: theme.typography.fontFamily.regular,
            },
          ]}
        />
        {rightIcon && <View style={styles.rightIcon}>{rightIcon}</View>}
      </View>
      {error && <Text style={[styles.errorText, { color: theme.colors.error, fontFamily: theme.typography.fontFamily.regular }]}>{error}</Text>}
    </View>
  );
};

// Progress Bar
export const ProgressBar = ({ progress = 0, color, style, height = 8 }) => {
  const { theme } = useThemeStore();
  const barColor = color || theme.colors.primary;

  return (
    <View
      style={[
        styles.progressContainer,
        {
          backgroundColor: theme.colors.surfaceVariant,
          height: height,
          borderRadius: height / 2,
        },
        style,
      ]}
    >
      <View
        style={[
          styles.progressBar,
          {
            backgroundColor: barColor,
            width: `${Math.min(100, Math.max(0, progress))}%`,
            borderRadius: height / 2,
          },
        ]}
      />
    </View>
  );
};

// Badge
export const Badge = ({ text, variant = 'neutral', style }) => {
  const { theme } = useThemeStore();

  const getBadgeColors = () => {
    switch (variant) {
      case 'success':
        return {
          bg: theme.colors.successBackground,
          text: theme.colors.success,
        };
      case 'warning':
        return {
          bg: theme.colors.warningBackground,
          text: theme.colors.warning,
        };
      case 'error':
        return {
          bg: theme.colors.errorBackground,
          text: theme.colors.error,
        };
      case 'info':
        return {
          bg: theme.colors.infoBackground,
          text: theme.colors.info,
        };
      default:
        return {
          bg: theme.colors.surfaceVariant,
          text: theme.colors.textMuted,
        };
    }
  };

  const colors = getBadgeColors();

  return (
    <View style={[styles.badge, { backgroundColor: colors.bg }, style]}>
      <Text
        style={[
          styles.badgeText,
          {
            color: colors.text,
            fontFamily: theme.typography.fontFamily.medium,
          },
        ]}
      >
        {text}
      </Text>
    </View>
  );
};

const styles = StyleSheet.create({
  button: {
    paddingVertical: 14,
    paddingHorizontal: 24,
    borderRadius: 12,
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 52,
    flexDirection: 'row',
  },
  buttonContent: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttonText: {
    fontSize: 16,
    fontWeight: '600',
  },
  disabled: {
    opacity: 0.6,
  },
  card: {
    borderRadius: 16,
    borderWidth: 1,
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.05,
    shadowRadius: 8,
    elevation: 3,
  },
  inputContainer: {
    marginBottom: 0,
  },
  label: {
    fontSize: 14,
    marginBottom: 8,
    marginLeft: 4,
  },
  inputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: 12,
    paddingHorizontal: 12,
    height: 52,
  },
  textInput: {
    flex: 1,
    fontSize: 16,
    height: '100%',
  },
  leftIcon: {
    marginRight: 10,
  },
  rightIcon: {
    marginLeft: 10,
  },
  errorText: {
    fontSize: 12,
    marginTop: 6,
    marginLeft: 4,
  },
  progressContainer: {
    overflow: 'hidden',
  },
  progressBar: {
    height: '100%',
  },
  badge: {
    paddingHorizontal: 10,
    paddingVertical: 6,
    borderRadius: 8,
    alignSelf: 'flex-start',
  },
  badgeText: {
    fontSize: 12,
    fontWeight: '600',
    textTransform: 'capitalize',
  },
});
