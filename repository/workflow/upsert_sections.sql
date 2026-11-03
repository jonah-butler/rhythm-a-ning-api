INSERT INTO workflow_sections (section_id, workflow_id, color, name, description)
SELECT t.section_id, $1::uuid, t.color, t.name, t.description
FROM UNNEST(
  $2::uuid[],    -- section ids
  $3::varchar[], -- colors
  $4::varchar[], -- names
  $5::text[]     -- descriptions
) AS t(section_id, color, name, description)
ON CONFLICT (section_id) DO UPDATE
SET
  color = EXCLUDED.color,
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  updated_at = NOW()
-- never take over a section that belongs to another workflow
WHERE workflow_sections.workflow_id = EXCLUDED.workflow_id
RETURNING section_id, workflow_id, color, name, description, created_at, updated_at;
