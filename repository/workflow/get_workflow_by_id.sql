SELECT
  workflow_id,
  user_id,
  name,
  description,
  created_at,
  updated_at
FROM workflows
WHERE workflow_id = $1
  AND user_id = $2;
