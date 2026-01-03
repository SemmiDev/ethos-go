-- ============================================================================
-- HABIT TRACKING APPLICATION - DOWN MIGRATION (ROLLBACK)
-- ============================================================================

-- Drop tables in dependency order
DROP TABLE IF EXISTS push_subscriptions CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS habit_vacations CASCADE;
DROP TABLE IF EXISTS habit_stats CASCADE;
DROP TABLE IF EXISTS habit_logs CASCADE;
DROP TABLE IF EXISTS habits CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop functions/triggers if any existed (legacy cleanup)
DROP FUNCTION IF EXISTS update_habit_streak CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
