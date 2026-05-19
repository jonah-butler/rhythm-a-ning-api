CREATE TABLE IF NOT EXISTS permissions (
	permission_id SERIAL PRIMARY KEY,
	name          TEXT NOT NULL UNIQUE
);