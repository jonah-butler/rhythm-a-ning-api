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