-- ===========================
-- user_tokens alterations   |
-- ===========================
ALTER TABLE user_tokens
	DROP CONSTRAINT IF EXISTS user_tokens_user_id_type_unique;
ALTER TABLE user_tokens
	ADD CONSTRAINT user_tokens_user_id_type_unique UNIQUE (user_id, type);