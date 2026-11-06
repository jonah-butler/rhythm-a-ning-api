SELECT
  section_id,
  workflow_id,
  color,
  name,
  description,
  created_at,
  updated_at
FROM workflow_sections
WHERE workflow_id = $1
ORDER BY created_at, name;
