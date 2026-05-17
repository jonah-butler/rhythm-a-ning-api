-- eventually merge this into a single sql command for
-- handling all token types
WITH inserted AS (
  INSERT INTO user_tokens (user_id, token, type, expires_at)
  VALUES ($1, $2, 2, $3) -- 2 = refresh_token type
  ON CONFLICT (user_id, type) DO UPDATE
  SET
    token = EXCLUDED.token,
    expires_at = EXCLUDED.expires_at
  RETURNING expires_at
)

SELECT EXISTS (SELECT 1 FROM inserted WHERE expires_at > NOW()) AS inserted;