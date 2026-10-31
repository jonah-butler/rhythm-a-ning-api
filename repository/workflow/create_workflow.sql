INSERT INTO
  workflows (
    user_id,
    name
  )
VALUES(
  $1, -- user_id
  $2  -- name
)
RETURNING
  workflow_id;