import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { useThemeStore } from '../stores/themeStore';

// Components
import CustomTabBar from './CustomTabBar';
import Header from '../components/Header';

// Screens
import DashboardScreen from '../screens/DashboardScreen';
import HabitsScreen from '../screens/HabitsScreen';
import HabitDetailScreen from '../screens/HabitDetailScreen';
import AnalyticsScreen from '../screens/AnalyticsScreen';
import SettingsScreen from '../screens/SettingsScreen';

const Tab = createBottomTabNavigator();
const HabitsStack = createNativeStackNavigator();

// Habits stack for nested navigation
function HabitsStackNavigator() {
  const { theme } = useThemeStore();

  return (
    <HabitsStack.Navigator
      screenOptions={{
        header: ({ navigation, route, options }) => {
          return <Header title={options.title || route.name} onBack={navigation.canGoBack() ? navigation.goBack : undefined} />;
        },
      }}
    >
      <HabitsStack.Screen name="HabitsList" component={HabitsScreen} options={{ title: 'My Habits' }} />
      <HabitsStack.Screen name="HabitDetail" component={HabitDetailScreen} options={{ title: 'Habit Details' }} />
    </HabitsStack.Navigator>
  );
}

export default function MainTabs() {
  return (
    <Tab.Navigator
      tabBar={(props) => <CustomTabBar {...props} />}
      screenOptions={{
        headerShown: true,
        header: ({ navigation, route, options }) => {
          return <Header title={route.name} />;
        },
      }}
    >
      <Tab.Screen
        name="Dashboard"
        component={DashboardScreen}
        options={{
          headerShown: false, // Dashboard has its own internal header
        }}
      />
      <Tab.Screen name="Habits" component={HabitsStackNavigator} options={{ headerShown: false }} />
      <Tab.Screen name="Analytics" component={AnalyticsScreen} options={{ title: 'Analytics' }} />
      <Tab.Screen name="Settings" component={SettingsScreen} options={{ title: 'Settings' }} />
    </Tab.Navigator>
  );
}
