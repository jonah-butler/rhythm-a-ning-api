SELECT workflow_id
FROM workflows
WHERE workflow_id = $1
  AND user_id = $2
FOR UPDATE;
