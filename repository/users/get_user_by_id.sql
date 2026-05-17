SELECT
  user_id,
  username,
  email,
  password
FROM users
WHERE user_id = $1;