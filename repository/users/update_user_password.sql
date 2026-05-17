WITH updated AS (
  UPDATE users
  SET password = $1
  WHERE user_id = $2
  RETURNING user_id
)

SELECT EXISTS (SELECT 1 FROM updated) AS updated;