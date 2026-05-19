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