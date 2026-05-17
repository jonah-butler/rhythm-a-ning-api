SELECT
  user_id
FROM users
WHERE username = $1
AND email = $2;