-- ===========================
-- users table               |
-- ===========================
CREATE TABLE IF NOT EXISTS users (
	user_id    SERIAL PRIMARY KEY,
	username   VARCHAR(50) NOT NULL UNIQUE,
	email      VARCHAR(100) NOT NULL UNIQUE,
	password   VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);