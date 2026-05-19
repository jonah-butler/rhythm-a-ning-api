-- permissions
INSERT INTO permissions (name)
VALUES
	('standard'),
	('subscriber'),
	('god')
ON CONFLICT (name)
DO NOTHING;