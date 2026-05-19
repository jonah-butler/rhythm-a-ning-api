-- ===========================
-- rhythms alterations       |
-- ===========================
ALTER TABLE rhythms
	ADD COLUMN IF NOT EXISTS level_id INT REFERENCES levels(level_id);