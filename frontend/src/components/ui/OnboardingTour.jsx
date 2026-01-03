import { useState, useEffect, useCallback } from 'react';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { X, ChevronRight, ChevronLeft, LayoutDashboard, Target, Command, Sun, BarChart3, Sparkles } from 'lucide-react';

/**
 * Onboarding Tour Component
 * Step-by-step guide for new users
 */

// Tour steps configuration
const TOUR_STEPS = [
  {
    id: 'welcome',
    title: 'Welcome to Ethos! 👋',
    description: "Let's take a quick tour to help you get started with building better habits.",
    icon: Sparkles,
    position: 'center',
  },
  {
    id: 'sidebar',
    title: 'Navigation Sidebar',
    description: 'Use the sidebar to navigate between Dashboard, Habits, Analytics, and Settings.',
    icon: LayoutDashboard,
    target: '[data-tour="sidebar"]',
    position: 'right',
  },
  {
    id: 'create-habit',
    title: 'Create Your First Habit',
    description: 'Click the "Add Habit" button to create a new habit you want to track.',
    icon: Target,
    target: '[data-tour="create-habit"]',
    position: 'bottom',
  },
  {
    id: 'command-palette',
    title: 'Quick Commands',
    description: 'Press ⌘K (or Ctrl+K) anytime to open the command palette for quick navigation.',
    icon: Command,
    position: 'center',
  },
  {
    id: 'theme',
    title: 'Dark Mode',
    description: 'Toggle between light and dark mode using the sun/moon icon in the header.',
    icon: Sun,
    target: '[data-tour="theme-toggle"]',
    position: 'bottom-left',
  },
  {
    id: 'analytics',
    title: 'Track Your Progress',
    description: 'Visit Analytics to see your completion rates, streaks, and progress over time.',
    icon: BarChart3,
    position: 'center',
  },
  {
    id: 'complete',
    title: "You're All Set! 🎉",
    description: "Start building better habits today. We'll celebrate your achievements along the way!",
    icon: Sparkles,
    position: 'center',
  },
];

// Tour store with persistence
export const useTourStore = create(
  persist(
    (set, get) => ({
      hasSeenTour: false,
      isActive: false,
      currentStep: 0,

      startTour: () => set({ isActive: true, currentStep: 0 }),

      endTour: () => {
        set({ isActive: false, hasSeenTour: true, currentStep: 0 });
      },

      nextStep: () => {
        const { currentStep } = get();
        if (currentStep < TOUR_STEPS.length - 1) {
          set({ currentStep: currentStep + 1 });
        } else {
          get().endTour();
        }
      },

      prevStep: () => {
        const { currentStep } = get();
        if (currentStep > 0) {
          set({ currentStep: currentStep - 1 });
        }
      },

      skipTour: () => {
        set({ isActive: false, hasSeenTour: true });
      },

      resetTour: () => {
        set({ hasSeenTour: false, isActive: false, currentStep: 0 });
      },
    }),
    {
      name: 'ethos-onboarding',
    }
  )
);

// Tooltip position calculator
const getTooltipPosition = (position) => {
  const positions = {
    center: 'fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2',
    top: 'absolute bottom-full left-1/2 -translate-x-1/2 mb-3',
    bottom: 'absolute top-full left-1/2 -translate-x-1/2 mt-3',
    left: 'absolute right-full top-1/2 -translate-y-1/2 mr-3',
    right: 'absolute left-full top-1/2 -translate-y-1/2 ml-3',
    'bottom-left': 'absolute top-full right-0 mt-3',
    'bottom-right': 'absolute top-full left-0 mt-3',
  };
  return positions[position] || positions.center;
};

// Progress dots component
const ProgressDots = ({ total, current }) => (
  <div className="flex gap-1.5 justify-center">
    {Array.from({ length: total }).map((_, i) => (
      <div
        key={i}
        className={`
                    w-2 h-2 rounded-full transition-all duration-200
                    ${i === current ? 'bg-primary w-6' : i < current ? 'bg-primary/60' : 'bg-base-300'}
                `}
      />
    ))}
  </div>
);

// Main Tour Component
export function OnboardingTour() {
  const { isActive, currentStep, nextStep, prevStep, skipTour, endTour } = useTourStore();
  const [targetRect, setTargetRect] = useState(null);

  const step = TOUR_STEPS[currentStep];
  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === TOUR_STEPS.length - 1;

  // Find and highlight target element
  useEffect(() => {
    if (!isActive || !step.target) {
      setTargetRect(null);
      return;
    }

    const targetElement = document.querySelector(step.target);
    if (targetElement) {
      const rect = targetElement.getBoundingClientRect();
      setTargetRect(rect);

      // Scroll element into view if needed
      targetElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }, [isActive, currentStep, step.target]);

  // Keyboard navigation
  useEffect(() => {
    if (!isActive) return;

    const handleKeyDown = (e) => {
      if (e.key === 'Escape') {
        skipTour();
      } else if (e.key === 'ArrowRight' || e.key === 'Enter') {
        nextStep();
      } else if (e.key === 'ArrowLeft' && !isFirstStep) {
        prevStep();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isActive, isFirstStep, nextStep, prevStep, skipTour]);

  if (!isActive) return null;

  const Icon = step.icon;

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 z-[90]">
        {/* Dark overlay with spotlight cutout */}
        <div
          className="absolute inset-0 bg-black/60 transition-all duration-300"
          style={
            targetRect
              ? {
                  clipPath: `polygon(
                            0 0, 100% 0, 100% 100%, 0 100%, 0 0,
                            ${targetRect.left - 8}px ${targetRect.top - 8}px,
                            ${targetRect.left - 8}px ${targetRect.bottom + 8}px,
                            ${targetRect.right + 8}px ${targetRect.bottom + 8}px,
                            ${targetRect.right + 8}px ${targetRect.top - 8}px,
                            ${targetRect.left - 8}px ${targetRect.top - 8}px
                        )`,
                }
              : undefined
          }
        />

        {/* Highlight ring around target */}
        {targetRect && (
          <div
            className="absolute border-2 border-primary rounded-lg animate-pulse pointer-events-none"
            style={{
              left: targetRect.left - 8,
              top: targetRect.top - 8,
              width: targetRect.width + 16,
              height: targetRect.height + 16,
            }}
          />
        )}
      </div>

      {/* Tooltip Card */}
      <div
        className={`
                    z-[95] w-full max-w-md p-6
                    bg-base-100 rounded-xl shadow-2xl border border-base-300
                    animate-in fade-in zoom-in-95 duration-300
                    ${step.position === 'center' ? 'fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2' : 'fixed'}
                `}
        style={
          step.position !== 'center' && targetRect
            ? {
                left: step.position.includes('right')
                  ? targetRect.right + 16
                  : step.position.includes('left')
                  ? targetRect.left - 420
                  : targetRect.left + targetRect.width / 2 - 200,
                top: step.position.includes('bottom')
                  ? targetRect.bottom + 16
                  : step.position.includes('top')
                  ? targetRect.top - 200
                  : targetRect.top + targetRect.height / 2 - 100,
              }
            : undefined
        }
      >
        {/* Close button */}
        <button onClick={skipTour} className="absolute top-4 right-4 p-1 text-base-content/40 hover:text-base-content transition-colors">
          <X size={20} />
        </button>

        {/* Icon */}
        <div className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-4">
          <Icon className="text-primary" size={24} />
        </div>

        {/* Content */}
        <h3 className="text-lg font-semibold text-base-content mb-2">{step.title}</h3>
        <p className="text-sm text-base-content/70 leading-relaxed mb-6">{step.description}</p>

        {/* Progress */}
        <ProgressDots total={TOUR_STEPS.length} current={currentStep} />

        {/* Actions */}
        <div className="flex items-center justify-between mt-6">
          <button onClick={skipTour} className="text-sm text-base-content/50 hover:text-base-content transition-colors">
            Skip tour
          </button>

          <div className="flex gap-2">
            {!isFirstStep && (
              <button
                onClick={prevStep}
                className="flex items-center gap-1 px-3 py-2 text-sm font-medium text-base-content/70 hover:text-base-content transition-colors"
              >
                <ChevronLeft size={16} />
                Back
              </button>
            )}
            <button
              onClick={isLastStep ? endTour : nextStep}
              className="flex items-center gap-1 px-4 py-2 bg-primary text-primary-content text-sm font-medium rounded-lg hover:bg-primary/90 transition-colors"
            >
              {isLastStep ? 'Get Started' : 'Next'}
              {!isLastStep && <ChevronRight size={16} />}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}

// Hook to auto-start tour for new users
export function useAutoStartTour() {
  const { hasSeenTour, startTour } = useTourStore();

  useEffect(() => {
    // Start tour after a short delay for new users
    if (!hasSeenTour) {
      const timer = setTimeout(() => {
        startTour();
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [hasSeenTour, startTour]);
}

// Help button to restart tour
export function TourRestartButton() {
  const { startTour } = useTourStore();

  return (
    <button onClick={startTour} className="text-sm text-primary hover:text-primary/80 font-medium transition-colors">
      Restart Tour
    </button>
  );
}
