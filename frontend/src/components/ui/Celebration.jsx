import { useState, useEffect, useCallback } from 'react';
import { create } from 'zustand';

/**
 * Confetti/Celebration Animation System
 * Triggers celebratory animations for achievements
 */

// Colors for confetti
const CONFETTI_COLORS = [
  '#3B82F6', // Blue
  '#10B981', // Green
  '#F59E0B', // Yellow
  '#EF4444', // Red
  '#8B5CF6', // Purple
  '#EC4899', // Pink
  '#06B6D4', // Cyan
  '#F97316', // Orange
];

// Celebration store for global access
export const useCelebrationStore = create((set) => ({
  isActive: false,
  type: 'confetti', // 'confetti' | 'fireworks' | 'stars'
  message: '',

  celebrate: (options = {}) => {
    set({
      isActive: true,
      type: options.type || 'confetti',
      message: options.message || '',
    });

    // Auto-hide after duration
    setTimeout(() => {
      set({ isActive: false, message: '' });
    }, options.duration || 3000);
  },

  stopCelebration: () => set({ isActive: false, message: '' }),
}));

// Individual confetti piece
const ConfettiPiece = ({ index, color }) => {
  const randomX = Math.random() * 100;
  const randomDelay = Math.random() * 0.5;
  const randomDuration = 2 + Math.random() * 2;
  const randomRotation = Math.random() * 360;
  const size = 8 + Math.random() * 8;
  const shape = Math.random() > 0.5 ? 'rounded-full' : 'rounded-sm';

  return (
    <div
      className={`absolute ${shape}`}
      style={{
        left: `${randomX}%`,
        top: '-20px',
        width: `${size}px`,
        height: `${size}px`,
        backgroundColor: color,
        transform: `rotate(${randomRotation}deg)`,
        animation: `confetti-fall ${randomDuration}s ease-out ${randomDelay}s forwards`,
        opacity: 0,
      }}
    />
  );
};

// Star piece for star celebration
const StarPiece = ({ index }) => {
  const randomX = 20 + Math.random() * 60;
  const randomY = 20 + Math.random() * 60;
  const randomDelay = Math.random() * 0.3;
  const size = 20 + Math.random() * 30;

  return (
    <div
      className="absolute text-yellow-400"
      style={{
        left: `${randomX}%`,
        top: `${randomY}%`,
        fontSize: `${size}px`,
        animation: `star-pop 0.6s ease-out ${randomDelay}s forwards`,
        opacity: 0,
      }}
    >
      ⭐
    </div>
  );
};

// Main Celebration Component
export function Celebration() {
  const { isActive, type, message } = useCelebrationStore();
  const [pieces, setPieces] = useState([]);

  useEffect(() => {
    if (isActive) {
      // Generate pieces based on type
      const count = type === 'stars' ? 12 : 50;
      const newPieces = Array.from({ length: count }).map((_, i) => ({
        id: i,
        color: CONFETTI_COLORS[Math.floor(Math.random() * CONFETTI_COLORS.length)],
      }));
      setPieces(newPieces);
    } else {
      setPieces([]);
    }
  }, [isActive, type]);

  if (!isActive) return null;

  return (
    <div className="fixed inset-0 pointer-events-none z-[100] overflow-hidden">
      {/* Confetti/Stars */}
      {type === 'confetti' && pieces.map((piece) => <ConfettiPiece key={piece.id} index={piece.id} color={piece.color} />)}

      {type === 'stars' && pieces.map((piece) => <StarPiece key={piece.id} index={piece.id} />)}

      {/* Center Message */}
      {message && (
        <div className="absolute inset-0 flex items-center justify-center">
          <div
            className="
                            px-8 py-4 bg-base-100/95 backdrop-blur-sm rounded-2xl
                            shadow-2xl border border-base-300
                            animate-in zoom-in-95 fade-in duration-300
                        "
          >
            <p className="text-2xl font-bold text-base-content text-center">{message}</p>
          </div>
        </div>
      )}

      {/* CSS Animations */}
      <style>{`
                @keyframes confetti-fall {
                    0% {
                        opacity: 1;
                        transform: translateY(0) rotate(0deg) scale(1);
                    }
                    100% {
                        opacity: 0;
                        transform: translateY(100vh) rotate(720deg) scale(0.5);
                    }
                }

                @keyframes star-pop {
                    0% {
                        opacity: 0;
                        transform: scale(0) rotate(0deg);
                    }
                    50% {
                        opacity: 1;
                        transform: scale(1.2) rotate(180deg);
                    }
                    100% {
                        opacity: 0;
                        transform: scale(1) rotate(360deg);
                    }
                }

                @keyframes firework-burst {
                    0% {
                        opacity: 1;
                        transform: scale(0);
                    }
                    50% {
                        opacity: 1;
                    }
                    100% {
                        opacity: 0;
                        transform: scale(2);
                    }
                }
            `}</style>
    </div>
  );
}

// Hook for easy celebration triggering
export function useCelebration() {
  const celebrate = useCelebrationStore((state) => state.celebrate);

  const triggerConfetti = useCallback(
    (message) => {
      celebrate({ type: 'confetti', message });
    },
    [celebrate]
  );

  const triggerStars = useCallback(
    (message) => {
      celebrate({ type: 'stars', message });
    },
    [celebrate]
  );

  const triggerHabitComplete = useCallback(() => {
    celebrate({
      type: 'confetti',
      message: '🎉 Habit Complete!',
      duration: 2500,
    });
  }, [celebrate]);

  const triggerStreak = useCallback(
    (days) => {
      celebrate({
        type: 'stars',
        message: `🔥 ${days} Day Streak!`,
        duration: 3500,
      });
    },
    [celebrate]
  );

  const triggerAllComplete = useCallback(() => {
    celebrate({
      type: 'confetti',
      message: '✨ All Habits Done!',
      duration: 4000,
    });
  }, [celebrate]);

  return {
    celebrate,
    triggerConfetti,
    triggerStars,
    triggerHabitComplete,
    triggerStreak,
    triggerAllComplete,
  };
}

// Milestone celebrations helper
export const MILESTONE_STREAKS = [7, 14, 21, 30, 50, 100, 365];

export function shouldCelebrateStreak(streak) {
  return MILESTONE_STREAKS.includes(streak);
}
