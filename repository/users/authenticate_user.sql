SELECT
  u.user_id
  ,u.username
  ,u.email
  ,u.password
FROM users u
EXISTS (
  SELECT 1
  FROM user_tokens ut
  WHERE ut.user_id = u.user_id
    AND ut.type = 1
) AS account_pending
WHERE email = $1;