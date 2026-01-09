// Spacing scale (matching Tailwind defaults)
export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  '2xl': 48,
};

// Typography
export const typography = {
  fontFamily: {
    regular: 'Inter_400Regular',
    medium: 'Inter_500Medium',
    semibold: 'Inter_600SemiBold',
    bold: 'Inter_700Bold',
  },
  fontSize: {
    xs: 12,
    sm: 14,
    base: 16,
    lg: 18,
    xl: 20,
    '2xl': 24,
    '3xl': 30,
  },
  lineHeight: {
    tight: 1.25,
    normal: 1.5,
    relaxed: 1.75,
  },
};

// Border radius
export const borderRadius = {
  sm: 4,
  md: 8,
  lg: 12,
  xl: 16,
  full: 9999,
};

// Shadows
export const shadows = {
  sm: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  md: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  lg: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.15,
    shadowRadius: 8,
    elevation: 5,
  },
};

export const lightTheme = {
  dark: false,
  colors: {
    // Primary: Deep Navy Blue - Trust & Security
    primary: '#0A2540',
    primaryContent: '#FFFFFF',

    // Secondary: Slate Blue - Professional accent
    secondary: '#1E3A5F',
    secondaryContent: '#FFFFFF',

    // Accent: Gold/Amber - Premium feel
    accent: '#C4A052',
    accentContent: '#0A2540',

    // Base: Clean whites and light grays
    background: '#F8FAFC',
    surface: '#FFFFFF',
    surfaceVariant: '#F1F5F9',
    border: '#E2E8F0',

    // Text
    text: '#1E293B',
    textMuted: 'rgba(30, 41, 59, 0.6)',
    textInverse: '#FFFFFF',

    // Semantic colors
    success: '#047857',
    successBackground: 'rgba(4, 120, 87, 0.1)',
    warning: '#B45309',
    warningBackground: 'rgba(180, 83, 9, 0.1)',
    error: '#B91C1C',
    errorBackground: 'rgba(185, 28, 28, 0.1)',
    info: '#0369A1',
    infoBackground: 'rgba(3, 105, 161, 0.1)',
  },
  typography,
  spacing,
  borderRadius,
  shadows,
};

export const darkTheme = {
  dark: true,
  colors: {
    // Primary: Lighter Blue for dark mode visibility
    primary: '#3B82F6',
    primaryContent: '#FFFFFF',

    // Secondary: Muted blue
    secondary: '#1E40AF',
    secondaryContent: '#FFFFFF',

    // Accent: Brighter Gold for dark mode
    accent: '#F59E0B',
    accentContent: '#0F172A',

    // Base: Dark backgrounds
    background: '#0F172A',
    surface: '#1E293B',
    surfaceVariant: '#334155',
    border: '#334155',

    // Text
    text: '#F1F5F9',
    textMuted: 'rgba(241, 245, 249, 0.6)',
    textInverse: '#0F172A',

    // Semantic colors - brighter for dark mode
    success: '#34D399',
    successBackground: 'rgba(52, 211, 153, 0.15)',
    warning: '#FBBF24',
    warningBackground: 'rgba(251, 191, 36, 0.15)',
    error: '#F87171',
    errorBackground: 'rgba(248, 113, 113, 0.15)',
    info: '#38BDF8',
    infoBackground: 'rgba(56, 189, 248, 0.15)',
  },
  typography,
  spacing,
  borderRadius,
  shadows,
};
