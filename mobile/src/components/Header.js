import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Platform } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { ChevronLeft } from 'lucide-react-native';
import { useThemeStore } from '../stores/themeStore';

export default function Header({ title, subtitle, onBack, rightAction, transparent = false, style }) {
  const { theme } = useThemeStore();
  const insets = useSafeAreaInsets();

  return (
    <View
      style={[
        styles.container,
        {
          backgroundColor: transparent ? 'transparent' : theme.colors.surface,
          paddingTop: insets.top,
          borderBottomWidth: transparent ? 0 : 1,
          borderBottomColor: theme.colors.border,
        },
        style,
      ]}
    >
      <View style={styles.content}>
        {/* Left Action (Back or Spacer) */}
        <View style={styles.left}>
          {onBack && (
            <TouchableOpacity onPress={onBack} style={styles.backButton}>
              <ChevronLeft color={theme.colors.text} size={28} />
            </TouchableOpacity>
          )}
        </View>

        {/* Center Title */}
        <View style={styles.center}>
          <Text style={[styles.title, { color: theme.colors.text }]}>{title}</Text>
          {subtitle && <Text style={[styles.subtitle, { color: theme.colors.textMuted }]}>{subtitle}</Text>}
        </View>

        {/* Right Action */}
        <View style={styles.right}>{rightAction}</View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    zIndex: 100,
  },
  content: {
    height: 56,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
  },
  left: {
    width: 48,
    alignItems: 'flex-start',
  },
  center: {
    flex: 1,
    alignItems: 'center',
  },
  right: {
    width: 48,
    alignItems: 'flex-end',
  },
  title: {
    fontSize: 17,
    fontWeight: '600',
    textAlign: 'center',
  },
  subtitle: {
    fontSize: 12,
    marginTop: 2,
    textAlign: 'center',
  },
  backButton: {
    padding: 8,
    marginLeft: -8,
  },
});
