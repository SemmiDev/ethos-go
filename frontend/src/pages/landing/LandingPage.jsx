import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Shield, Target, BarChart3, Bell, Zap, ArrowRight, Star, Moon, Lock } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { LanguageToggle } from '../../components/ui/LanguageToggle';

// Animation variants
const fadeInUp = {
  hidden: { opacity: 0, y: 30 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.6, ease: 'easeOut' },
  },
};

const fadeIn = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { duration: 0.6 },
  },
};

const staggerContainer = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.2,
    },
  },
};

const scaleIn = {
  hidden: { opacity: 0, scale: 0.8 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.5, ease: 'easeOut' },
  },
};

// Feature data
const features = [
  {
    icon: Target,
    title: 'Track Any Habit',
    description: 'Daily, weekly, or custom schedules. Track habits your way with flexible goal setting.',
  },
  {
    icon: BarChart3,
    title: 'Visual Analytics',
    description: 'Beautiful charts and insights to understand your patterns and celebrate progress.',
  },
  {
    icon: Bell,
    title: 'Smart Reminders',
    description: 'Never miss a habit with customizable notifications that keep you on track.',
  },
  {
    icon: Zap,
    title: 'Streak Tracking',
    description: 'Build momentum with streak tracking and milestone celebrations.',
  },
  {
    icon: Moon,
    title: 'Dark Mode',
    description: 'Easy on your eyes with beautiful light and dark themes.',
  },
  {
    icon: Lock,
    title: 'Secure & Private',
    description: 'Your data is encrypted and never shared. Your habits, your privacy.',
  },
];

// How it works steps
const steps = [
  {
    step: '01',
    title: 'Create Your Habits',
    description: 'Add habits you want to build or break. Set your goals and schedule.',
  },
  {
    step: '02',
    title: 'Track Daily',
    description: 'Log your progress each day with a simple tap. It takes seconds.',
  },
  {
    step: '03',
    title: 'Build Streaks',
    description: 'Watch your streaks grow and celebrate milestones along the way.',
  },
  {
    step: '04',
    title: 'See Results',
    description: 'Visualize your progress and become the person you want to be.',
  },
];

// Stats
const stats = [
  { value: '10K+', label: 'Active Users' },
  { value: '1M+', label: 'Habits Tracked' },
  { value: '98%', label: 'Success Rate' },
  { value: '4.9★', label: 'App Rating' },
];

// Testimonials
const testimonials = [
  {
    quote: "Ethos helped me build a daily meditation habit. After 90 days, I can't imagine my mornings without it.",
    author: 'Sarah K.',
    role: 'Product Designer',
    avatar: 'S',
  },
  {
    quote: "The streak feature is incredibly motivating. I've never stuck with a habit program this long before.",
    author: 'Michael R.',
    role: 'Software Engineer',
    avatar: 'M',
  },
  {
    quote: 'Clean, simple, effective. Exactly what I needed to finally build consistent workout habits.',
    author: 'Emily T.',
    role: 'Entrepreneur',
    avatar: 'E',
  },
];

export function LandingPage() {
  const { t } = useTranslation();

  return (
    <div className="min-h-screen bg-base-100 overflow-hidden">
      {/* Navigation */}
      <motion.nav
        className="sticky top-0 z-50 bg-base-100/80 backdrop-blur-md border-b border-base-200"
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        transition={{ duration: 0.5, ease: 'easeOut' }}
      >
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2">
              <motion.img
                src="/logo.jpg"
                alt="Ethos Logo"
                className="w-12 h-12 rounded-lg object-cover"
                whileHover={{ scale: 1.05, rotate: 5 }}
                whileTap={{ scale: 0.95 }}
              />
            </Link>

            {/* Nav Links */}
            <div className="hidden md:flex items-center gap-8">
              {['Features', 'How it Works', 'Testimonials'].map((item, i) => (
                <motion.a
                  key={item}
                  href={`#${item.toLowerCase().replace(/\s+/g, '-')}`}
                  className="text-sm font-medium text-base-content/70 hover:text-primary transition-colors"
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.1 * i }}
                >
                  {item}
                </motion.a>
              ))}
            </div>

            {/* Auth Buttons & Language */}
            <div className="flex items-center gap-3">
              <LanguageToggle size="sm" />
              <Link to="/login" className="text-sm font-medium text-base-content/70 hover:text-base-content transition-colors">
                {t('nav.login')}
              </Link>
              <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                <Link to="/register" className="px-4 py-2 bg-primary text-primary-content text-sm font-medium rounded-lg hover:bg-primary/90 transition-colors">
                  {t('landing.hero.cta')}
                </Link>
              </motion.div>
            </div>
          </div>
        </div>
      </motion.nav>

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Animated Background */}
        <motion.div
          className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-secondary/5"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1 }}
        />
        <motion.div
          className="absolute top-0 right-0 w-96 h-96 bg-primary/10 rounded-full blur-3xl"
          animate={{
            x: [0, 30, 0],
            y: [0, -20, 0],
          }}
          transition={{ duration: 8, repeat: Infinity, ease: 'easeInOut' }}
        />
        <motion.div
          className="absolute bottom-0 left-0 w-96 h-96 bg-secondary/10 rounded-full blur-3xl"
          animate={{
            x: [0, -30, 0],
            y: [0, 20, 0],
          }}
          transition={{ duration: 10, repeat: Infinity, ease: 'easeInOut' }}
        />

        <div className="relative max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-24 lg:py-32">
          <div className="text-center max-w-3xl mx-auto">
            {/* Badge */}
            <motion.div
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary/10 rounded-full mb-8"
              variants={scaleIn}
              initial="hidden"
              animate="visible"
            >
              <motion.div animate={{ rotate: [0, 10, -10, 0] }} transition={{ duration: 2, repeat: Infinity }}>
                <Zap size={14} className="text-primary" />
              </motion.div>
              <span className="text-sm font-medium text-primary">Build better habits, one day at a time</span>
            </motion.div>

            {/* Headline */}
            <motion.h1
              className="text-4xl sm:text-5xl lg:text-6xl font-bold text-base-content leading-tight mb-6"
              variants={fadeInUp}
              initial="hidden"
              animate="visible"
            >
              The Simple Way to{' '}
              <motion.span
                className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-secondary"
                animate={{
                  backgroundPosition: ['0%', '100%', '0%'],
                }}
                transition={{ duration: 5, repeat: Infinity }}
              >
                Transform Your Life
              </motion.span>
            </motion.h1>

            {/* Subheadline */}
            <motion.p
              className="text-lg sm:text-xl text-base-content/60 mb-10 max-w-2xl mx-auto leading-relaxed"
              variants={fadeInUp}
              initial="hidden"
              animate="visible"
              transition={{ delay: 0.2 }}
            >
              Ethos helps you build positive habits and break bad ones with streak tracking, smart reminders, and beautiful analytics.
            </motion.p>

            {/* CTA Buttons */}
            <motion.div
              className="flex flex-col sm:flex-row items-center justify-center gap-4"
              variants={fadeInUp}
              initial="hidden"
              animate="visible"
              transition={{ delay: 0.3 }}
            >
              <motion.div whileHover={{ scale: 1.03 }} whileTap={{ scale: 0.97 }}>
                <Link
                  to="/register"
                  className="w-full sm:w-auto flex items-center justify-center gap-2 px-8 py-4 bg-primary text-primary-content font-semibold rounded-xl hover:bg-primary/90 shadow-lg shadow-primary/20 transition-all hover:shadow-xl hover:shadow-primary/30"
                >
                  Start Free Today
                  <motion.div animate={{ x: [0, 5, 0] }} transition={{ duration: 1.5, repeat: Infinity }}>
                    <ArrowRight size={18} />
                  </motion.div>
                </Link>
              </motion.div>
              <motion.a
                href="#how-it-works"
                className="w-full sm:w-auto flex items-center justify-center gap-2 px-8 py-4 bg-base-200 text-base-content font-semibold rounded-xl hover:bg-base-300 transition-colors"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                See How It Works
              </motion.a>
            </motion.div>

            {/* Social Proof */}
            <motion.div
              className="flex items-center justify-center gap-6 mt-12 pt-12 border-t border-base-200"
              variants={fadeIn}
              initial="hidden"
              animate="visible"
              transition={{ delay: 0.5 }}
            >
              <div className="flex -space-x-2">
                {['A', 'B', 'C', 'D', 'E'].map((letter, i) => (
                  <motion.div
                    key={i}
                    className="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-sm font-medium border-2 border-base-100"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: 0.6 + i * 0.1 }}
                    whileHover={{ scale: 1.1, zIndex: 10 }}
                  >
                    {letter}
                  </motion.div>
                ))}
              </div>
              <div className="text-left">
                <div className="flex items-center gap-1">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <motion.div key={i} initial={{ opacity: 0, scale: 0 }} animate={{ opacity: 1, scale: 1 }} transition={{ delay: 0.8 + i * 0.05 }}>
                      <Star size={16} className="fill-yellow-400 text-yellow-400" />
                    </motion.div>
                  ))}
                </div>
                <p className="text-sm text-base-content/60">
                  <span className="font-semibold text-base-content">10,000+</span> people building habits
                </p>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-16 bg-base-200/50 border-y border-base-200">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            className="grid grid-cols-2 lg:grid-cols-4 gap-8"
            variants={staggerContainer}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-100px' }}
          >
            {stats.map((stat, i) => (
              <motion.div key={i} className="text-center" variants={scaleIn}>
                <motion.p
                  className="text-3xl sm:text-4xl font-bold text-primary mb-1"
                  initial={{ opacity: 0 }}
                  whileInView={{ opacity: 1 }}
                  viewport={{ once: true }}
                >
                  {stat.value}
                </motion.p>
                <p className="text-sm text-base-content/60">{stat.label}</p>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-24">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Section Header */}
          <motion.div
            className="text-center max-w-2xl mx-auto mb-16"
            variants={fadeInUp}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-100px' }}
          >
            <p className="text-primary font-semibold mb-3">Features</p>
            <h2 className="text-3xl sm:text-4xl font-bold text-base-content mb-4">Everything you need to build lasting habits</h2>
            <p className="text-base-content/60">Simple yet powerful tools designed to help you succeed.</p>
          </motion.div>

          {/* Features Grid */}
          <motion.div
            className="grid sm:grid-cols-2 lg:grid-cols-3 gap-8"
            variants={staggerContainer}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-50px' }}
          >
            {features.map((feature, i) => (
              <motion.div
                key={i}
                className="p-6 bg-base-100 border border-base-200 rounded-2xl hover:border-primary/30 hover:shadow-lg transition-all duration-300 group cursor-pointer"
                variants={fadeInUp}
                whileHover={{ y: -5, transition: { duration: 0.2 } }}
              >
                <motion.div
                  className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-4 group-hover:bg-primary/20 transition-colors"
                  whileHover={{ rotate: 5, scale: 1.1 }}
                >
                  <feature.icon className="text-primary" size={24} />
                </motion.div>
                <h3 className="text-lg font-semibold text-base-content mb-2">{feature.title}</h3>
                <p className="text-base-content/60 text-sm leading-relaxed">{feature.description}</p>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* How It Works Section */}
      <section id="how-it-works" className="py-24 bg-base-200/30">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Section Header */}
          <motion.div
            className="text-center max-w-2xl mx-auto mb-16"
            variants={fadeInUp}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-100px' }}
          >
            <p className="text-primary font-semibold mb-3">How It Works</p>
            <h2 className="text-3xl sm:text-4xl font-bold text-base-content mb-4">Start building habits in minutes</h2>
            <p className="text-base-content/60">Four simple steps to transform your daily routine.</p>
          </motion.div>

          {/* Steps */}
          <motion.div
            className="grid sm:grid-cols-2 lg:grid-cols-4 gap-8"
            variants={staggerContainer}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-50px' }}
          >
            {steps.map((step, i) => (
              <motion.div key={i} className="relative" variants={fadeInUp}>
                {/* Connector line */}
                {i < steps.length - 1 && (
                  <motion.div
                    className="hidden lg:block absolute top-8 left-full w-full h-0.5 bg-gradient-to-r from-primary/40 to-transparent"
                    initial={{ scaleX: 0 }}
                    whileInView={{ scaleX: 1 }}
                    viewport={{ once: true }}
                    transition={{ delay: 0.5 + i * 0.2 }}
                  />
                )}
                <div className="text-center">
                  <motion.div
                    className="w-16 h-16 rounded-2xl bg-primary text-primary-content text-2xl font-bold flex items-center justify-center mx-auto mb-4"
                    whileHover={{ scale: 1.1, rotate: 5 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    {step.step}
                  </motion.div>
                  <h3 className="text-lg font-semibold text-base-content mb-2">{step.title}</h3>
                  <p className="text-base-content/60 text-sm">{step.description}</p>
                </div>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* Testimonials Section */}
      <section id="testimonials" className="py-24">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Section Header */}
          <motion.div
            className="text-center max-w-2xl mx-auto mb-16"
            variants={fadeInUp}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-100px' }}
          >
            <p className="text-primary font-semibold mb-3">Testimonials</p>
            <h2 className="text-3xl sm:text-4xl font-bold text-base-content mb-4">Loved by thousands of habit builders</h2>
            <p className="text-base-content/60">See what our community has to say about their journey.</p>
          </motion.div>

          {/* Testimonials Grid */}
          <motion.div
            className="grid md:grid-cols-3 gap-8"
            variants={staggerContainer}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: '-50px' }}
          >
            {testimonials.map((testimonial, i) => (
              <motion.div key={i} className="p-6 bg-base-100 border border-base-200 rounded-2xl" variants={fadeInUp} whileHover={{ y: -5 }}>
                {/* Stars */}
                <div className="flex gap-1 mb-4">
                  {[1, 2, 3, 4, 5].map((j) => (
                    <Star key={j} size={16} className="fill-yellow-400 text-yellow-400" />
                  ))}
                </div>
                {/* Quote */}
                <p className="text-base-content/80 mb-6 leading-relaxed">"{testimonial.quote}"</p>
                {/* Author */}
                <div className="flex items-center gap-3">
                  <motion.div
                    className="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white font-medium"
                    whileHover={{ scale: 1.1 }}
                  >
                    {testimonial.avatar}
                  </motion.div>
                  <div>
                    <p className="font-semibold text-base-content">{testimonial.author}</p>
                    <p className="text-sm text-base-content/50">{testimonial.role}</p>
                  </div>
                </div>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-primary to-secondary p-12 text-center"
            initial={{ opacity: 0, y: 50 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
          >
            {/* Background decoration */}
            <motion.div
              className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full blur-3xl"
              animate={{
                scale: [1, 1.2, 1],
                opacity: [0.1, 0.2, 0.1],
              }}
              transition={{ duration: 4, repeat: Infinity }}
            />
            <motion.div
              className="absolute bottom-0 left-0 w-64 h-64 bg-white/10 rounded-full blur-3xl"
              animate={{
                scale: [1, 1.3, 1],
                opacity: [0.1, 0.15, 0.1],
              }}
              transition={{ duration: 5, repeat: Infinity }}
            />

            <div className="relative z-10">
              <motion.h2
                className="text-3xl sm:text-4xl font-bold text-white mb-4"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: 0.2 }}
              >
                Ready to build better habits?
              </motion.h2>
              <motion.p
                className="text-white/80 mb-8 max-w-lg mx-auto"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: 0.3 }}
              >
                Join thousands of people who are already transforming their lives with Ethos. Start free today.
              </motion.p>
              <motion.div
                className="flex flex-col sm:flex-row items-center justify-center gap-4"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: 0.4 }}
              >
                <motion.div whileHover={{ scale: 1.03 }} whileTap={{ scale: 0.97 }}>
                  <Link
                    to="/register"
                    className="w-full sm:w-auto flex items-center justify-center gap-2 px-8 py-4 bg-white text-primary font-semibold rounded-xl hover:bg-white/90 transition-colors"
                  >
                    Get Started Free
                    <ArrowRight size={18} />
                  </Link>
                </motion.div>
                <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                  <Link
                    to="/login"
                    className="w-full sm:w-auto flex items-center justify-center gap-2 px-8 py-4 bg-white/10 text-white font-semibold rounded-xl hover:bg-white/20 border border-white/20 transition-colors"
                  >
                    Sign In
                  </Link>
                </motion.div>
              </motion.div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 border-t border-base-200">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            className="flex flex-col md:flex-row items-center justify-between gap-6"
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
          >
            {/* Logo */}
            <motion.div className="flex items-center gap-2" whileHover={{ scale: 1.02 }}>
              <img src="/logo.jpg" alt="Ethos Logo" className="w-10 h-10 rounded-lg object-cover" />
            </motion.div>

            {/* Links */}
            <div className="flex items-center gap-8 text-sm text-base-content/60">
              <a href="#features" className="hover:text-base-content transition-colors">
                Features
              </a>
              <a href="#how-it-works" className="hover:text-base-content transition-colors">
                How it Works
              </a>
              <Link to="/help" className="hover:text-base-content transition-colors">
                Help
              </Link>
              <a href="#" className="hover:text-base-content transition-colors">
                Privacy
              </a>
              <a href="#" className="hover:text-base-content transition-colors">
                Terms
              </a>
            </div>

            {/* Copyright */}
            <p className="text-sm text-base-content/40">© 2026 Ethos. All rights reserved.</p>
          </motion.div>
        </div>
      </footer>
    </div>
  );
}
