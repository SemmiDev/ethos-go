import { useState, useEffect } from 'react';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { Download, X, Smartphone, Monitor, CheckCircle } from 'lucide-react';

/**
 * PWA Install Prompt Component
 * Shows install banner for installable PWAs
 */

// PWA store for install state
export const usePWAStore = create(
  persist(
    (set, get) => ({
      deferredPrompt: null,
      isInstallable: false,
      isInstalled: false,
      dismissedUntil: null, // Timestamp when user dismissed

      setDeferredPrompt: (prompt) => set({ deferredPrompt: prompt, isInstallable: !!prompt }),
      setInstalled: () => set({ isInstalled: true, isInstallable: false }),

      dismissPrompt: () => {
        // Dismiss for 7 days
        const dismissedUntil = Date.now() + 7 * 24 * 60 * 60 * 1000;
        set({ dismissedUntil });
      },

      shouldShowPrompt: () => {
        const { isInstallable, isInstalled, dismissedUntil } = get();
        if (isInstalled || !isInstallable) return false;
        if (dismissedUntil && Date.now() < dismissedUntil) return false;
        return true;
      },
    }),
    {
      name: 'ethos-pwa',
      partialize: (state) => ({
        isInstalled: state.isInstalled,
        dismissedUntil: state.dismissedUntil,
      }),
    }
  )
);

// Register service worker and handle install prompt
export function usePWA() {
  const { setDeferredPrompt, setInstalled, deferredPrompt } = usePWAStore();

  useEffect(() => {
    // Register service worker
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker
        .register('/sw.js')
        .then((registration) => {
          console.log('[PWA] Service worker registered:', registration);
        })
        .catch((error) => {
          console.error('[PWA] Service worker registration failed:', error);
        });
    }

    // Listen for beforeinstallprompt
    const handleBeforeInstall = (e) => {
      e.preventDefault();
      console.log('[PWA] Install prompt captured');
      setDeferredPrompt(e);
    };

    // Listen for app installed
    const handleAppInstalled = () => {
      console.log('[PWA] App installed');
      setInstalled();
      setDeferredPrompt(null);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstall);
    window.addEventListener('appinstalled', handleAppInstalled);

    // Check if already installed (standalone mode)
    if (window.matchMedia('(display-mode: standalone)').matches) {
      setInstalled();
    }

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstall);
      window.removeEventListener('appinstalled', handleAppInstalled);
    };
  }, [setDeferredPrompt, setInstalled]);

  const install = async () => {
    if (!deferredPrompt) return false;

    deferredPrompt.prompt();
    const { outcome } = await deferredPrompt.userChoice;

    console.log('[PWA] Install outcome:', outcome);

    if (outcome === 'accepted') {
      setInstalled();
    }

    setDeferredPrompt(null);
    return outcome === 'accepted';
  };

  return { install };
}

// Mobile vs Desktop detection
function useDeviceType() {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkDevice = () => {
      setIsMobile(/Android|iPhone|iPad|iPod/i.test(navigator.userAgent));
    };
    checkDevice();
    window.addEventListener('resize', checkDevice);
    return () => window.removeEventListener('resize', checkDevice);
  }, []);

  return isMobile;
}

// Install Banner Component
export function PWAInstallBanner() {
  const { shouldShowPrompt, dismissPrompt, isInstalled } = usePWAStore();
  const { install } = usePWA();
  const isMobile = useDeviceType();
  const [isInstalling, setIsInstalling] = useState(false);

  if (!shouldShowPrompt() || isInstalled) return null;

  const handleInstall = async () => {
    setIsInstalling(true);
    await install();
    setIsInstalling(false);
  };

  return (
    <div className="fixed bottom-4 left-4 right-4 md:left-auto md:right-4 md:w-96 z-50 animate-in slide-in-from-bottom fade-in duration-300">
      <div className="bg-base-100 border border-base-300 rounded-xl shadow-2xl p-4">
        <div className="flex items-start gap-4">
          {/* Icon */}
          <div className="p-3 bg-primary/10 rounded-xl shrink-0">
            {isMobile ? <Smartphone className="text-primary" size={24} /> : <Monitor className="text-primary" size={24} />}
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <h4 className="font-semibold text-base-content mb-1">Install Ethos App</h4>
            <p className="text-sm text-base-content/60">Install for quick access, offline support, and a native app experience.</p>

            {/* Actions */}
            <div className="flex gap-2 mt-3">
              <button
                onClick={handleInstall}
                disabled={isInstalling}
                className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-content text-sm font-medium rounded-lg hover:bg-primary/90 disabled:opacity-50 transition-colors"
              >
                <Download size={16} />
                {isInstalling ? 'Installing...' : 'Install'}
              </button>
              <button onClick={dismissPrompt} className="px-4 py-2 text-sm font-medium text-base-content/60 hover:text-base-content transition-colors">
                Not now
              </button>
            </div>
          </div>

          {/* Close */}
          <button onClick={dismissPrompt} className="p-1 text-base-content/40 hover:text-base-content transition-colors">
            <X size={18} />
          </button>
        </div>
      </div>
    </div>
  );
}

// Compact install button for settings page
export function PWAInstallButton() {
  const { isInstallable, isInstalled } = usePWAStore();
  const { install } = usePWA();
  const [isInstalling, setIsInstalling] = useState(false);

  if (isInstalled) {
    return (
      <div className="flex items-center gap-2 text-success text-sm">
        <CheckCircle size={16} />
        App Installed
      </div>
    );
  }

  if (!isInstallable) {
    return null;
  }

  const handleInstall = async () => {
    setIsInstalling(true);
    await install();
    setIsInstalling(false);
  };

  return (
    <button
      onClick={handleInstall}
      disabled={isInstalling}
      className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-content text-sm font-medium rounded-lg hover:bg-primary/90 disabled:opacity-50 transition-colors"
    >
      <Download size={16} />
      {isInstalling ? 'Installing...' : 'Install App'}
    </button>
  );
}
