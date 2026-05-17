SELECT
  user_id,
  username,
  email,
  password
FROM users
WHERE email = $1;