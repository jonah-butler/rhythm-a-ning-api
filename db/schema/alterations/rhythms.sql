-- ===========================
-- rhythms alterations       |
-- ===========================
ALTER TABLE rhythms
	ADD COLUMN IF NOT EXISTS level_id INT REFERENCES levels(level_id);

ALTER TABLE rhythms
	ADD COLUMN IF NOT EXISTS name VARCHAR(250) NOT NULL,
	ADD COLUMN IF NOT EXISTS description TEXT;
