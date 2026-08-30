CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ----------------------
-- TABLE CREATION      
-- ----------------------                    
-- ++++++++++++++++++++++                     
-- ++++++++++++++++++++++                      

-- ================================
-- lookup table for subdivisions  |
-- ================================
CREATE TABLE IF NOT EXISTS subdivision_types (
	subdivision_id SERIAL PRIMARY KEY,
	name           VARCHAR(50) NOT NULL UNIQUE -- 'base', 'duplet', 'triplet', etc
);

-- ================================
-- lookup table for levels        |
-- ================================
CREATE TABLE IF NOT EXISTS levels (
	level_id       SERIAL PRIMARY KEY,
	name           VARCHAR(50) NOT NULL UNIQUE -- 'beginner', 'intermediate', 'advanced'
);

CREATE TABLE IF NOT EXISTS permissions (
	permission_id SERIAL PRIMARY KEY,
	name          TEXT NOT NULL UNIQUE
);

-- ===========================
-- users table               |
-- ===========================
CREATE TABLE IF NOT EXISTS users (
	user_id    		UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username   		VARCHAR(50) NOT NULL UNIQUE,
	email      		VARCHAR(100) NOT NULL UNIQUE,
	password   		VARCHAR(255) NOT NULL,
	permission_id INT REFERENCES permissions(permission_id),
	created_at 		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at 		TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =======================================
-- tags for labeling workflows/rhythms   |
-- =======================================
CREATE TABLE IF NOT EXISTS tags (
	tag_id      SERIAL PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE,
	description TEXT,
	user_id UUID REFERENCES users(user_id)
);


-- =======================
-- primary rhythms table |
-- =======================
CREATE TABLE IF NOT EXISTS rhythms (
	rhythm_id 			 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	bpm              SMALLINT NOT NULL,
	beats            SMALLINT NOT NULL,
	sounds           JSONB, -- 
	subdivision      INT NOT NULL REFERENCES subdivision_types(subdivision_id),
	state            SMALLINT[] NOT NULL,
	is_poly          BOOLEAN NOT NULL DEFAULT FALSE,
	poly_beats       SMALLINT, -- can be NULL
	poly_subdivision INT REFERENCES subdivision_types(subdivision_id), -- can be NULL
	poly_state       SMALLINT[],
	poly_sounds      JSONB, -- 
	user_id          UUID REFERENCES users(user_id), -- NULL is global
	level_id 				 INT REFERENCES levels(level_id),
	name 						 VARCHAR(250) NOT NULL,
	description 		 TEXT, -- can be null
	tag_id           INT REFERENCES tags(tag_id),
	created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 

-- =================
-- workflows table |
-- =================
CREATE TABLE IF NOT EXISTS workflows (
	workflow_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id 	  UUID, -- null is global
	name        VARCHAR(250) NOT NULL,
	description VARCHAR(500),
	created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ====================
-- workflow sections |
-- ===================
CREATE TABLE IF NOT EXISTS workflow_sections (
	section_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	color VARCHAR(7), -- hex can be null
	name VARCHAR(250) NOT NULL,
	description TEXT
);

-- ===================================================
-- table for rhythms stored within a workflow        |
-- references its workflow, and rhythm               |
-- with the addition of position within the workflow |
-- ===================================================
CREATE TABLE IF NOT EXISTS workflow_rhythms (
    workflow_rhythm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id        UUID NOT NULL REFERENCES workflows(workflow_id) ON DELETE CASCADE,
    rhythm_id          UUID NOT NULL REFERENCES rhythms(rhythm_id) ON DELETE CASCADE,
		section_id         UUID REFERENCES workflow_sections(section_id) ON DELETE SET NULL,
    measures           SMALLINT NOT NULL,
    position           SMALLINT NOT NULL
);

-- =============================
-- user token types            |
-- =============================
CREATE TABLE IF NOT EXISTS user_token_types (
	token_type_id SERIAL PRIMARY KEY,
	name          VARCHAR(100) NOT NULL UNIQUE
);

-- =============================
-- user tokens                 |
-- =============================
CREATE TABLE IF NOT EXISTS user_tokens (
    user_token_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token         TEXT NOT NULL UNIQUE,
		type          INT NOT NULL REFERENCES user_token_types(token_type_id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		expires_at    TIMESTAMPTZ NOT NULL
);

-- ===================================
-- associating workflows with tags  |
-- ==================================
CREATE TABLE IF NOT EXISTS workflow_tags (
    workflow_id UUID NOT NULL REFERENCES workflows(workflow_id) ON DELETE CASCADE,
    tag_id      INT NOT NULL REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (workflow_id, tag_id)
);

-- ===================================
-- associating rhythms with tags  |
-- ==================================
CREATE TABLE IF NOT EXISTS rhythm_tags (
    rhythm_id UUID NOT NULL REFERENCES rhythms(rhythm_id) ON DELETE CASCADE,
    tag_id    INT NOT NULL REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (rhythm_id, tag_id)
);


-- ----------------------
-- TABLE ALTERATIONS      
-- ----------------------                    
-- ++++++++++++++++++++++                     
-- ++++++++++++++++++++++ 

ALTER TABLE user_tokens
	DROP CONSTRAINT IF EXISTS user_tokens_user_id_type_unique;
ALTER TABLE user_tokens
	ADD CONSTRAINT user_tokens_user_id_type_unique UNIQUE (user_id, type);

-- ----------------------
-- DEFAULT INSERTIONS      
-- ----------------------                    
-- ++++++++++++++++++++++                     
-- ++++++++++++++++++++++ 

-- user_token_types
INSERT INTO user_token_types (name)
VALUES
	('account verification'),
	('refresh token'),
	('password reset')
ON CONFLICT (name)
DO NOTHING;

-- permissions
INSERT INTO permissions (name)
VALUES
	('standard'),
	('subscriber'),
	('god')
ON CONFLICT (name)
DO NOTHING;

-- levels
INSERT INTO levels (name)
VALUES
	('beginner'),
	('intermediate'),
	('advanced')
ON CONFLICT (name)
DO NOTHING;

-- subdivision_types
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

-- tags
INSERT INTO tags (name, description)
VALUES
  ('Fundamentals', 'Exercises to gain familiarity with the metronome and develop a strong rhythmic foundation'),
  ('Odd Meters', 'Explore odd meters to develop a strong internal clock'),
  ('Polyrhythms', ''),
  ('African Rhythms', ''),
  ('Speed Development', ''),
  ('Slow Motion', 'To speed up sometimes requires slowing down'),
  ('Displacement', ''),
	('Training', ''),
	('World Rhythms', '')
	ON CONFLICT (name)
	DO NOTHING;