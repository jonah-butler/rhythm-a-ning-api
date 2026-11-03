UPDATE workflows
SET
  name = $1,
  description = $2,
  updated_at = NOW()
WHERE
  workflow_id = $3
  AND user_id = $4
RETURNING
  workflow_id,
  user_id,
  name,
  description,
  created_at,
  updated_at;