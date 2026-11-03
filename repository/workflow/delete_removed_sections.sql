-- $2 must be an empty array (not NULL) when no sections are kept
DELETE FROM workflow_sections
WHERE workflow_id = $1
  AND section_id <> ALL($2::uuid[]);
