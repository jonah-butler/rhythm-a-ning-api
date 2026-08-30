WITH updated AS (
  UPDATE user_tokens
  SET 
    token = $1,
    expires_at = $2
  WHERE
    user_id = $3
  AND
    type = 1 -- 1 = user registration type
  RETURNING expires_at
)

SELECT EXISTS (SELECT 1 FROM updated WHERE expires_at > NOW()) AS updated;