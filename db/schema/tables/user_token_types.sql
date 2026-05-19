-- =============================
-- user token types            |
-- =============================
CREATE TABLE IF NOT EXISTS user_token_types (
	token_type_id SERIAL PRIMARY KEY,
	name          VARCHAR(100) NOT NULL UNIQUE
);