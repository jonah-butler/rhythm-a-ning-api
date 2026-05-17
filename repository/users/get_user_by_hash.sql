SELECT
  user_id,
  expires_at
FROM user_tokens
WHERE token = $1
AND type = $2;