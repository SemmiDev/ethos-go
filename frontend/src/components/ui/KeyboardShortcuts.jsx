import { useState, useEffect, useCallback } from 'react';
import {
  X,
  Command,
  Keyboard,
  LayoutDashboard,
  Target,
  BarChart3,
  Settings,
  HelpCircle,
  Plus,
  Sun,
  LogOut,
  Search,
  ArrowUp,
  ArrowDown,
  CornerDownLeft,
} from 'lucide-react';

/**
 * Keyboard Shortcuts Help Modal
 * Press ? to show all available shortcuts
 */

// Shortcut categories and their shortcuts
const SHORTCUT_CATEGORIES = [
  {
    title: 'Navigation',
    shortcuts: [
      { keys: ['G', 'D'], description: 'Go to Dashboard', icon: LayoutDashboard },
      { keys: ['G', 'H'], description: 'Go to Habits', icon: Target },
      { keys: ['G', 'A'], description: 'Go to Analytics', icon: BarChart3 },
      { keys: ['G', 'S'], description: 'Go to Settings', icon: Settings },
      { keys: ['G', '?'], description: 'Go to Help', icon: HelpCircle },
    ],
  },
  {
    title: 'Actions',
    shortcuts: [
      { keys: ['⌘', 'K'], description: 'Open Command Palette', icon: Search },
      { keys: ['C'], description: 'Create New Habit', icon: Plus },
      { keys: ['T'], description: 'Toggle Theme', icon: Sun },
      { keys: ['?'], description: 'Show Keyboard Shortcuts', icon: Keyboard },
    ],
  },
  {
    title: 'Command Palette',
    shortcuts: [
      { keys: ['↑', '↓'], description: 'Navigate Items', icon: ArrowUp },
      { keys: ['↵'], description: 'Select Item', icon: CornerDownLeft },
      { keys: ['Esc'], description: 'Close', icon: X },
    ],
  },
  {
    title: 'General',
    shortcuts: [{ keys: ['Esc'], description: 'Close Modal / Cancel', icon: X }],
  },
];

// Keyboard key component
const KeyboardKey = ({ children, isSpecial = false }) => (
  <kbd
    className={`
            inline-flex items-center justify-center min-w-[24px] h-6 px-1.5
            text-xs font-medium rounded
            ${isSpecial ? 'bg-primary/10 text-primary border border-primary/20' : 'bg-base-200 text-base-content border border-base-300'}
            shadow-sm
        `}
  >
    {children}
  </kbd>
);

// Single shortcut row
const ShortcutRow = ({ keys, description, icon: Icon }) => (
  <div className="flex items-center justify-between py-2.5 border-b border-base-200 last:border-0">
    <div className="flex items-center gap-3">
      <Icon size={16} className="text-base-content/50" />
      <span className="text-sm text-base-content">{description}</span>
    </div>
    <div className="flex items-center gap-1">
      {keys.map((key, i) => (
        <span key={i} className="flex items-center">
          <KeyboardKey isSpecial={key === '⌘'}>{key}</KeyboardKey>
          {i < keys.length - 1 && <span className="text-base-content/30 text-xs mx-0.5">then</span>}
        </span>
      ))}
    </div>
  </div>
);

// Main Keyboard Shortcuts Modal
export function KeyboardShortcutsModal({ isOpen, onClose }) {
  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 animate-in fade-in duration-150" onClick={onClose} />

      {/* Modal */}
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4 pointer-events-none">
        <div
          className="
                        w-full max-w-2xl max-h-[80vh] overflow-hidden
                        bg-base-100 rounded-xl shadow-2xl border border-base-300
                        pointer-events-auto
                        animate-in fade-in zoom-in-95 duration-200
                    "
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="flex items-center justify-between px-6 py-4 border-b border-base-200">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <Keyboard className="text-primary" size={20} />
              </div>
              <div>
                <h2 className="text-lg font-semibold text-base-content">Keyboard Shortcuts</h2>
                <p className="text-sm text-base-content/50">Navigate faster with these shortcuts</p>
              </div>
            </div>
            <button onClick={onClose} className="p-2 text-base-content/40 hover:text-base-content hover:bg-base-200 rounded-lg transition-colors">
              <X size={20} />
            </button>
          </div>

          {/* Content */}
          <div className="p-6 overflow-y-auto max-h-[calc(80vh-100px)]">
            <div className="grid gap-6 md:grid-cols-2">
              {SHORTCUT_CATEGORIES.map((category) => (
                <div key={category.title}>
                  <h3 className="text-xs font-semibold text-base-content/40 uppercase tracking-wider mb-3">{category.title}</h3>
                  <div className="bg-base-50 rounded-lg px-4">
                    {category.shortcuts.map((shortcut, i) => (
                      <ShortcutRow key={i} {...shortcut} />
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Footer */}
          <div className="px-6 py-3 border-t border-base-200 bg-base-200/30">
            <p className="text-xs text-base-content/40 text-center">
              Press <KeyboardKey>?</KeyboardKey> anytime to show this help
            </p>
          </div>
        </div>
      </div>
    </>
  );
}

// Hook for keyboard shortcuts
export function useKeyboardShortcuts() {
  const [isHelpOpen, setIsHelpOpen] = useState(false);
  const [pendingKey, setPendingKey] = useState(null);

  const handleShortcut = useCallback((action) => {
    // Import navigate dynamically to avoid hook issues
    import('react-router-dom').then(({ useNavigate }) => {
      // This won't work directly, we need to use an event system
    });
  }, []);

  useEffect(() => {
    let timeout = null;

    const handleKeyDown = (e) => {
      // Skip if user is typing in an input
      const target = e.target;
      if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
        return;
      }

      // Show help with ?
      if (e.key === '?' && !e.ctrlKey && !e.metaKey) {
        e.preventDefault();
        setIsHelpOpen(true);
        return;
      }

      // Close help with Escape
      if (e.key === 'Escape' && isHelpOpen) {
        e.preventDefault();
        setIsHelpOpen(false);
        return;
      }

      // Handle G + key navigation
      if (pendingKey === 'g') {
        e.preventDefault();
        const routes = {
          d: '/dashboard',
          h: '/habits',
          a: '/analytics',
          s: '/settings',
          '?': '/help',
        };
        const route = routes[e.key.toLowerCase()];
        if (route) {
          window.location.href = route;
        }
        setPendingKey(null);
        clearTimeout(timeout);
        return;
      }

      // Start G sequence
      if (e.key.toLowerCase() === 'g' && !e.ctrlKey && !e.metaKey) {
        setPendingKey('g');
        timeout = setTimeout(() => setPendingKey(null), 1000);
        return;
      }

      // Quick actions
      if (!e.ctrlKey && !e.metaKey) {
        switch (e.key.toLowerCase()) {
          case 'c':
            // Trigger create habit (dispatch custom event)
            window.dispatchEvent(new CustomEvent('shortcut:create-habit'));
            break;
          case 't':
            // Trigger theme toggle
            window.dispatchEvent(new CustomEvent('shortcut:toggle-theme'));
            break;
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      if (timeout) clearTimeout(timeout);
    };
  }, [isHelpOpen, pendingKey]);

  return {
    isHelpOpen,
    openHelp: () => setIsHelpOpen(true),
    closeHelp: () => setIsHelpOpen(false),
  };
}

// Shortcut indicator component (shows pending keys)
export function ShortcutIndicator() {
  const [pending, setPending] = useState('');

  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key.toLowerCase() === 'g') {
        setPending('g');
        setTimeout(() => setPending(''), 1000);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  if (!pending) return null;

  return (
    <div className="fixed bottom-4 right-4 z-40 animate-in fade-in slide-in-from-bottom duration-200">
      <div className="flex items-center gap-2 px-4 py-2 bg-base-100 border border-base-300 rounded-lg shadow-lg">
        <KeyboardKey>{pending}</KeyboardKey>
        <span className="text-sm text-base-content/50">+ ?</span>
      </div>
    </div>
  );
}
