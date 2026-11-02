INSERT INTO
  workflows (
    user_id,
    name,
    description
  )
VALUES(
  $1, -- user_id
  $2, -- name
  $3  -- description
)
RETURNING
  workflow_id,
  user_id,
  name,
  description,
  created_at,
  updated_at;
