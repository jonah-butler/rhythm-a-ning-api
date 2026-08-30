SELECT
  token
FROM
  user_tokens
WHERE
  user_id = $1
AND
  type = $2;