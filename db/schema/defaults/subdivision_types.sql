-- subdivision_types
INSERT INTO subdivision_types (name)
VALUES
	('base'),
	('duplet'),
	('triplet'),
	('quadruplet'),
	('quintuplet'),
	('sextuplet'),
	('septuplet'),
	('octuplet'),
	('nonuplet'),
	('decuplet')
ON CONFLICT (name)
DO NOTHING;