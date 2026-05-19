-- ===========================
-- users alterations         |
-- ===========================
ALTER TABLE users
	ADD COLUMN IF NOT EXISTS permission_id INT REFERENCES permissions(permission_id);