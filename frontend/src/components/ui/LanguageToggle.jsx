import { useLanguageStore } from '../../stores/languageStore';
import { languages } from '../../i18n';
import { Globe } from 'lucide-react';

export function LanguageToggle({ showLabel = false, size = 'md' }) {
  const { language, setLanguage } = useLanguageStore();

  const currentLang = languages.find((l) => l.code === language) || languages[0];
  const nextLang = languages.find((l) => l.code !== language) || languages[1];

  const handleToggle = () => {
    setLanguage(nextLang.code);
  };

  const sizeClasses = {
    sm: 'w-8 h-8 text-sm',
    md: 'w-10 h-10 text-base',
    lg: 'w-12 h-12 text-lg',
  };

  return (
    <button
      onClick={handleToggle}
      className={`
        ${sizeClasses[size]}
        flex items-center justify-center
        rounded-lg
        bg-base-200 hover:bg-base-300
        text-base-content
        transition-all duration-200
        hover:scale-105 active:scale-95
        border border-base-300
        shadow-sm
      `}
      title={`Switch to ${nextLang.name}`}
      aria-label={`Switch language to ${nextLang.name}`}
    >
      <span className="text-lg">{currentLang.flag}</span>
      {showLabel && <span className="ml-2 text-sm font-medium">{currentLang.code.toUpperCase()}</span>}
    </button>
  );
}

// Dropdown version for settings page
export function LanguageSelector() {
  const { language, setLanguage } = useLanguageStore();

  return (
    <div className="flex items-center gap-3">
      <Globe className="w-5 h-5 text-base-content/50" />
      <select
        value={language}
        onChange={(e) => setLanguage(e.target.value)}
        className="flex-1 px-3 py-2.5 text-sm bg-base-100 border border-base-300 rounded-md text-base-content focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20"
      >
        {languages.map((lang) => (
          <option key={lang.code} value={lang.code}>
            {lang.flag} {lang.name}
          </option>
        ))}
      </select>
    </div>
  );
}

// Animated toggle with flags (like dark mode toggle)
export function LanguageToggleAnimated() {
  const { language, setLanguage } = useLanguageStore();

  const isEnglish = language === 'en';

  return (
    <button
      onClick={() => setLanguage(isEnglish ? 'id' : 'en')}
      className="relative flex items-center w-16 h-8 p-1 rounded-full bg-base-200 border border-base-300 transition-all duration-300 hover:bg-base-300"
      aria-label="Toggle language"
    >
      {/* Sliding indicator */}
      <div
        className={`
          absolute w-6 h-6 rounded-full bg-primary shadow-md
          transition-all duration-300 ease-in-out
          flex items-center justify-center text-sm
          ${isEnglish ? 'left-1' : 'left-[calc(100%-1.75rem)]'}
        `}
      >
        <span className="text-primary-content text-xs font-bold">{isEnglish ? 'EN' : 'ID'}</span>
      </div>

      {/* Flag icons */}
      <div className="flex w-full justify-between px-1.5">
        <span className={`text-sm transition-opacity duration-200 ${isEnglish ? 'opacity-0' : 'opacity-50'}`}>🇺🇸</span>
        <span className={`text-sm transition-opacity duration-200 ${isEnglish ? 'opacity-50' : 'opacity-0'}`}>🇮🇩</span>
      </div>
    </button>
  );
}
