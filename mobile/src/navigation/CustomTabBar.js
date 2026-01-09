import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Platform } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { LayoutDashboard, CheckSquare, BarChart2, Settings } from 'lucide-react-native';
import { useThemeStore } from '../stores/themeStore';

const TabIcon = ({ name, color, size }) => {
  switch (name) {
    case 'Dashboard':
      return <LayoutDashboard color={color} size={size} />;
    case 'Habits':
      return <CheckSquare color={color} size={size} />;
    case 'Analytics':
      return <BarChart2 color={color} size={size} />;
    case 'Settings':
      return <Settings color={color} size={size} />;
    default:
      return <LayoutDashboard color={color} size={size} />;
  }
};

export default function CustomTabBar({ state, descriptors, navigation }) {
  const { theme } = useThemeStore();
  const insets = useSafeAreaInsets();

  return (
    <View
      style={[
        styles.container,
        {
          backgroundColor: theme.colors.surface,
          borderTopColor: theme.colors.border,
          paddingBottom: Math.max(insets.bottom, 20),
        },
      ]}
    >
      {state.routes.map((route, index) => {
        const { options } = descriptors[route.key];
        const label = options.tabBarLabel !== undefined ? options.tabBarLabel : options.title !== undefined ? options.title : route.name;

        const isFocused = state.index === index;

        const onPress = () => {
          const event = navigation.emit({
            type: 'tabPress',
            target: route.key,
            canPreventDefault: true,
          });

          if (!isFocused && !event.defaultPrevented) {
            navigation.navigate(route.name);
          }
        };

        const color = isFocused ? theme.colors.primary : theme.colors.textMuted;

        return (
          <TouchableOpacity
            key={route.key}
            accessibilityRole="button"
            accessibilityState={isFocused ? { selected: true } : {}}
            accessibilityLabel={options.tabBarAccessibilityLabel}
            testID={options.tabBarTestID}
            onPress={onPress}
            style={styles.tab}
          >
            <View style={[styles.iconContainer, isFocused && { backgroundColor: theme.colors.primary + '15' }]}>
              <TabIcon name={route.name} color={color} size={24} />
            </View>
            <Text style={[styles.label, { color }]}>{label}</Text>
          </TouchableOpacity>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    borderTopWidth: 1,
    paddingTop: 12,
    elevation: 8,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: -2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  tab: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  iconContainer: {
    padding: 8,
    borderRadius: 12,
    marginBottom: 4,
  },
  label: {
    fontSize: 11,
    fontWeight: '500',
  },
});
