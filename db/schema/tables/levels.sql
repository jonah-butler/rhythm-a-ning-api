-- ================================
-- lookup table for levels        |
-- ================================
CREATE TABLE IF NOT EXISTS levels (
	level_id       SERIAL PRIMARY KEY,
	name           VARCHAR(50) NOT NULL UNIQUE -- 'beginner', 'intermediate', 'advanced'
);