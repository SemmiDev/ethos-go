import daisyui from 'daisyui';

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [daisyui],
  daisyui: {
    themes: [
      {
        banking: {
          // Primary: Deep Navy Blue - Trust & Security
          primary: '#0A2540',
          'primary-content': '#FFFFFF',

          // Secondary: Slate Blue - Professional accent
          secondary: '#1E3A5F',
          'secondary-content': '#FFFFFF',

          // Accent: Gold/Amber - Premium feel
          accent: '#C4A052',
          'accent-content': '#0A2540',

          // Neutral: Dark grays
          neutral: '#1F2937',
          'neutral-content': '#F9FAFB',

          // Base: Clean whites and light grays
          'base-100': '#FFFFFF',
          'base-200': '#F8FAFC',
          'base-300': '#E2E8F0',
          'base-content': '#1E293B',

          // Info: Professional blue
          info: '#0369A1',
          'info-content': '#FFFFFF',

          // Success: Muted green
          success: '#047857',
          'success-content': '#FFFFFF',

          // Warning: Muted amber
          warning: '#B45309',
          'warning-content': '#FFFFFF',

          // Error: Muted red
          error: '#B91C1C',
          'error-content': '#FFFFFF',

          // Border radius - More conservative/professional
          '--rounded-box': '0.5rem',
          '--rounded-btn': '0.375rem',
          '--rounded-badge': '0.25rem',

          // Animations - Subtle
          '--animation-btn': '0.15s',
          '--animation-input': '0.15s',

          // Button focus scale - Minimal
          '--btn-focus-scale': '0.99',

          // Border width
          '--border-btn': '1px',
        },
      },
      {
        'banking-dark': {
          // Primary: Lighter Blue for dark mode visibility
          primary: '#3B82F6',
          'primary-content': '#FFFFFF',

          // Secondary: Muted blue
          secondary: '#1E40AF',
          'secondary-content': '#FFFFFF',

          // Accent: Brighter Gold for dark mode
          accent: '#F59E0B',
          'accent-content': '#0F172A',

          // Neutral: Light grays for dark mode
          neutral: '#374151',
          'neutral-content': '#F9FAFB',

          // Base: Dark backgrounds
          'base-100': '#0F172A',
          'base-200': '#1E293B',
          'base-300': '#334155',
          'base-content': '#F1F5F9',

          // Info: Brighter blue
          info: '#38BDF8',
          'info-content': '#0F172A',

          // Success: Brighter green
          success: '#34D399',
          'success-content': '#0F172A',

          // Warning: Brighter amber
          warning: '#FBBF24',
          'warning-content': '#0F172A',

          // Error: Brighter red
          error: '#F87171',
          'error-content': '#0F172A',

          // Border radius - Same as light
          '--rounded-box': '0.5rem',
          '--rounded-btn': '0.375rem',
          '--rounded-badge': '0.25rem',

          // Animations - Subtle
          '--animation-btn': '0.15s',
          '--animation-input': '0.15s',

          // Button focus scale - Minimal
          '--btn-focus-scale': '0.99',

          // Border width
          '--border-btn': '1px',
        },
      },
    ],
  },
};
