WITH inserted AS (
  INSERT INTO user_tokens (user_id, token, expires_at, type)
  VALUES($1, $2, $3, 1) -- 1 = user_registration type
  ON CONFLICT (user_id, type) DO UPDATE
  SET
    token = EXCLUDED.token,
    expires_at = EXCLUDED.expires_at
  RETURNING expires_at
)

SELECT EXISTS (SELECT 1 FROM inserted WHERE expires_at > NOW()) AS inserted;