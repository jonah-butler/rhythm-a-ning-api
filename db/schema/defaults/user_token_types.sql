-- user_token_types
INSERT INTO user_token_types (name)
VALUES
	('account verification'),
	('refresh token'),
	('password reset')
ON CONFLICT (name)
DO NOTHING;