import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { HelpCircle, MessageCircle, Mail, Book, ChevronDown, ExternalLink, Target, BarChart3, Bell, Shield, Zap, Heart } from 'lucide-react';
import { Header } from '../../components/layout/Sidebar';

// FAQ Item Component
const FAQItem = ({ question, answer }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="border-b border-base-200 last:border-0">
      <button onClick={() => setIsOpen(!isOpen)} className="w-full flex items-center justify-between py-4 text-left group">
        <span className="text-sm font-medium text-base-content group-hover:text-primary transition-colors">{question}</span>
        <ChevronDown
          size={18}
          className={`
                        text-base-content/40 shrink-0 ml-4
                        transition-transform duration-200
                        ${isOpen ? 'rotate-180' : ''}
                    `}
        />
      </button>
      <div
        className={`
                    overflow-hidden transition-all duration-300 ease-out
                    ${isOpen ? 'max-h-96 pb-4' : 'max-h-0'}
                `}
      >
        <p className="text-sm text-base-content/70 leading-relaxed pl-0">{answer}</p>
      </div>
    </div>
  );
};

// FAQ Category Component
const FAQCategory = ({ title, icon: Icon, questions }) => {
  const [isExpanded, setIsExpanded] = useState(true);

  return (
    <div className="bg-base-100 border border-base-300 rounded-lg overflow-hidden">
      <button onClick={() => setIsExpanded(!isExpanded)} className="w-full flex items-center gap-3 p-4 bg-base-200/50 hover:bg-base-200 transition-colors">
        <div className="p-2 bg-primary/10 rounded-lg">
          <Icon size={18} className="text-primary" />
        </div>
        <span className="text-base font-semibold text-base-content flex-1 text-left">{title}</span>
        <ChevronDown
          size={18}
          className={`
                        text-base-content/40
                        transition-transform duration-200
                        ${isExpanded ? 'rotate-180' : ''}
                    `}
        />
      </button>
      <div
        className={`
                    transition-all duration-300 ease-out
                    ${isExpanded ? 'max-h-[2000px]' : 'max-h-0 overflow-hidden'}
                `}
      >
        <div className="px-4">
          {questions.map((item, index) => (
            <FAQItem key={index} question={item.q} answer={item.a} />
          ))}
        </div>
      </div>
    </div>
  );
};

// Quick Link Card Component
const QuickLinkCard = ({ icon: Icon, title, description, href, external = false }) => (
  <a
    href={href}
    target={external ? '_blank' : undefined}
    rel={external ? 'noopener noreferrer' : undefined}
    className="block p-5 bg-base-100 border border-base-300 rounded-lg hover:border-primary/30 hover:shadow-md transition-all duration-200 group"
  >
    <div className="flex items-start gap-4">
      <div className="p-2.5 bg-primary/10 rounded-lg group-hover:bg-primary/20 transition-colors">
        <Icon size={20} className="text-primary" />
      </div>
      <div className="flex-1">
        <div className="flex items-center gap-2">
          <h3 className="font-semibold text-base-content group-hover:text-primary transition-colors">{title}</h3>
          {external && <ExternalLink size={14} className="text-base-content/40" />}
        </div>
        <p className="text-sm text-base-content/60 mt-1">{description}</p>
      </div>
    </div>
  </a>
);

export function HelpPage() {
  const { t } = useTranslation();

  // FAQ Data with translations
  const faqCategories = [
    {
      title: t('help.faq.gettingStarted.title'),
      icon: Zap,
      questions: [
        { q: t('help.faq.gettingStarted.q1'), a: t('help.faq.gettingStarted.a1') },
        { q: t('help.faq.gettingStarted.q2'), a: t('help.faq.gettingStarted.a2') },
      ],
    },
    {
      title: t('help.faq.tracking.title'),
      icon: Target,
      questions: [
        { q: t('help.faq.tracking.q1'), a: t('help.faq.tracking.a1') },
        { q: t('help.faq.tracking.q2'), a: t('help.faq.tracking.a2') },
      ],
    },
    {
      title: t('help.faq.analytics.title'),
      icon: BarChart3,
      questions: [
        { q: t('help.faq.analytics.q1'), a: t('help.faq.analytics.a1') },
        { q: t('help.faq.analytics.q2'), a: t('help.faq.analytics.a2') },
      ],
    },
    {
      title: t('help.faq.account.title'),
      icon: Shield,
      questions: [
        { q: t('help.faq.account.q1'), a: t('help.faq.account.a1') },
        { q: t('help.faq.account.q2'), a: t('help.faq.account.a2') },
      ],
    },
    {
      title: t('help.faq.notifications.title'),
      icon: Bell,
      questions: [
        { q: t('help.faq.notifications.q1'), a: t('help.faq.notifications.a1') },
        { q: t('help.faq.notifications.q2'), a: t('help.faq.notifications.a2') },
      ],
    },
  ];

  return (
    <div className="space-y-8">
      <Header title={t('help.title')} subtitle={t('help.subtitle')} />

      {/* Hero Section */}
      <div className="bg-gradient-to-br from-primary to-secondary rounded-xl p-8 text-white relative overflow-hidden">
        {/* Background Pattern */}
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white rounded-full blur-3xl" />
          <div className="absolute bottom-0 left-0 w-48 h-48 bg-white rounded-full blur-3xl" />
        </div>

        <div className="relative z-10 max-w-2xl">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2.5 bg-white/20 rounded-lg backdrop-blur">
              <HelpCircle size={24} />
            </div>
            <h2 className="text-2xl font-bold">{t('help.contact.title')}</h2>
          </div>
          <p className="text-white/80 leading-relaxed">{t('help.contact.subtitle')}</p>
        </div>
      </div>

      {/* Quick Links */}
      <div>
        <h2 className="text-lg font-semibold text-base-content mb-4">{t('help.quickLinks.title')}</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <QuickLinkCard icon={Book} title={t('help.quickLinks.dashboard')} description={t('help.contact.docsDesc')} href="/dashboard" />
          <QuickLinkCard
            icon={MessageCircle}
            title={t('help.contact.email')}
            description={t('help.contact.emailDesc')}
            href="mailto:support@ethos-app.com"
            external
          />
          <QuickLinkCard icon={Heart} title={t('help.contact.docs')} description={t('help.contact.docsDesc')} href="mailto:feedback@ethos-app.com" external />
        </div>
      </div>

      {/* FAQ Section */}
      <div>
        <h2 className="text-lg font-semibold text-base-content mb-4">{t('help.faq.title')}</h2>
        <div className="space-y-4">
          {faqCategories.map((category, index) => (
            <FAQCategory key={index} title={category.title} icon={category.icon} questions={category.questions} />
          ))}
        </div>
      </div>

      {/* Contact Section */}
      <div className="bg-base-100 border border-base-300 rounded-lg p-6">
        <div className="flex flex-col md:flex-row md:items-center gap-6">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-base-content mb-2">{t('help.contact.title')}</h3>
            <p className="text-sm text-base-content/60">{t('help.contact.subtitle')}</p>
          </div>
          <div className="flex flex-col sm:flex-row gap-3">
            <a
              href="mailto:support@ethos-app.com"
              className="inline-flex items-center justify-center gap-2 px-5 py-2.5 bg-primary text-primary-content rounded-lg font-medium text-sm hover:bg-primary/90 transition-colors"
            >
              <Mail size={16} />
              {t('help.contact.email')}
            </a>
          </div>
        </div>
      </div>

      {/* App Info */}
      <div className="text-center py-6 border-t border-base-200">
        <p className="text-sm text-base-content/50">Ethos Habit Tracker • Version 1.0.0</p>
        <p className="text-xs text-base-content/40 mt-1">Made with ❤️ for building better habits</p>
      </div>
    </div>
  );
}
