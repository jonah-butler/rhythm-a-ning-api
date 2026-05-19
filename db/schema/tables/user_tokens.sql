-- =============================
-- user tokens                 |
-- =============================
CREATE TABLE IF NOT EXISTS user_tokens (
    user_token_id SERIAL PRIMARY KEY,
    user_id       INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token         TEXT NOT NULL UNIQUE,
		type          INT NOT NULL REFERENCES user_token_types(token_type_id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		expires_at    TIMESTAMPTZ NOT NULL
);