-- ===========================
-- user_tokens index         |
-- ===========================
DROP INDEX IF EXISTS user_refresh_tokens_user_id_idx;    -- double check
DROP INDEX IF EXISTS user_refresh_tokens_expires_at_idx; -- double check
DROP TABLE IF EXISTS user_refresh_tokens;