-- ================================
-- lookup table for subdivisions  |
-- ================================
CREATE TABLE IF NOT EXISTS subdivision_types (
	subdivision_id SERIAL PRIMARY KEY,
	name           VARCHAR(50) NOT NULL UNIQUE -- 'base', 'duplet', 'triplet', etc
);

-- ===============================================================
-- populate subdivision lookup table with supported subdivisions |
-- ===============================================================
INSERT INTO subdivision_types (name)
VALUES
	('base'),
	('duplet'),
	('triplet'),
	('quadruplet'),
	('quintuplet'),
	('sextuplet'),
	('septuplet'),
	('octuplet'),
	('nonuplet'),
	('decuplet')
ON CONFLICT (name)
DO NOTHING;

-- =======================
-- primary rhythms table |
-- =======================
CREATE TABLE IF NOT EXISTS rhythms (
	rhythm_id        SERIAL PRIMARY KEY,
	bpm              SMALLINT NOT NULL,
	beats            SMALLINT NOT NULL,
	subdivision      INT NOT NULL REFERENCES subdivision_types(subdivision_id),
	state            SMALLINT[] NOT NULL,
	is_poly          BOOLEAN NOT NULL DEFAULT FALSE,
	poly_beats       SMALLINT, -- can be NULL
	poly_subdivision INT REFERENCES subdivision_types(subdivision_id), -- can be NULL
	poly_state       SMALLINT[],
	user_id          INT, -- NULL is global
	created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 

-- =================
-- workflows table |
-- =================
CREATE TABLE IF NOT EXISTS workflows (
	workflow_id SERIAL PRIMARY KEY,
	name        VARCHAR(250) NOT NULL,
	description VARCHAR(500),
	created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===================================================
-- table for rhythms stored within a workflow        |
-- references its workflow, and rhythm               |
-- with the addition of position within the workflow |
-- ===================================================
CREATE TABLE IF NOT EXISTS workflow_rhythms (
    workflow_rhythm_id SERIAL PRIMARY KEY,
    workflow_id        INT NOT NULL REFERENCES workflows(workflow_id) ON DELETE CASCADE,
    rhythm_id          INT NOT NULL REFERENCES rhythms(rhythm_id) ON DELETE CASCADE,
    measures           SMALLINT NOT NULL,
    position           SMALLINT NOT NULL
);

-- ================================
-- lookup table for levels        |
-- ================================
CREATE TABLE IF NOT EXISTS levels (
	level_id SERIAL PRIMARY KEY,
	name           VARCHAR(50) NOT NULL UNIQUE -- 'beginner', 'intermediate', 'advanced'
);

INSERT INTO levels (name)
VALUES
	('beginner'),
	('intermediate'),
	('advanced')
ON CONFLICT (name)
DO NOTHING;

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

ALTER TABLE users
	ADD COLUMN IF NOT EXISTS permission_id INT REFERENCES permissions(permission_id);

DROP INDEX IF EXISTS user_refresh_tokens_user_id_idx; -- double check
DROP INDEX IF EXISTS user_refresh_tokens_expires_at_idx; -- double check
DROP TABLE IF EXISTS user_refresh_tokens;

CREATE TABLE IF NOT EXISTS permissions (
	permission_id SERIAL PRIMARY KEY
	name TEXT NOT NULL UNIQUE
)

INSERT INTO user_permissions (name)
VALUES
	('standard'),
	('subscriber'),
	('god')
ON CONFLICT (name)
DO NOTHING;


-- =============================
-- user token types            |
-- =============================
CREATE TABLE IF NOT EXISTS user_token_types (
	token_type_id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO user_token_types (name)
VALUES
	('account verification'),
	('refresh token'),
	('password reset')
ON CONFLICT (name)
DO NOTHING;

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

ALTER TABLE user_tokens
	DROP CONSTRAINT IF EXISTS user_tokens_user_id_type_unique;
ALTER TABLE user_tokens
	ADD CONSTRAINT user_tokens_user_id_type_unique UNIQUE (user_id, type);

-- ===========================
-- rhythms table alterations |
-- ===========================
ALTER TABLE rhythms
ADD COLUMN IF NOT EXISTS level_id INT REFERENCES levels(level_id);

-- =============================
-- workflows table alterations |
-- =============================
ALTER TABLE workflows
	ADD COLUMN IF NOT EXISTS user_id INT; -- null is global