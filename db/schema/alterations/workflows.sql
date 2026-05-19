-- =============================
-- workflows alterations       |
-- =============================
ALTER TABLE workflows
	ADD COLUMN IF NOT EXISTS user_id INT; -- null is global