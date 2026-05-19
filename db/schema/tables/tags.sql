-- ===============================
-- tags for labeling workflows   |
-- ===============================
CREATE TABLE IF NOT EXISTS tags (
	tag_id      SERIAL PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE,
	description TEXT
);