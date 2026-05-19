-- ================================
-- lookup table for subdivisions  |
-- ================================
CREATE TABLE IF NOT EXISTS subdivision_types (
	subdivision_id SERIAL PRIMARY KEY,
	name           VARCHAR(50) NOT NULL UNIQUE -- 'base', 'duplet', 'triplet', etc
);