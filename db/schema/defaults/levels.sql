-- levels
INSERT INTO levels (name)
VALUES
	('beginner'),
	('intermediate'),
	('advanced')
ON CONFLICT (name)
DO NOTHING;