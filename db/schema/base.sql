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
-- rhythms table alterations |
-- ===========================
ALTER TABLE rhythms
ADD COLUMN IF NOT EXISTS level_id INT REFERENCES levels(level_id);