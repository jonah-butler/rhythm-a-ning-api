-- workflow_sections and any remaining workflow_rhythms rows cascade
DELETE FROM workflows
WHERE workflow_id = $1;
