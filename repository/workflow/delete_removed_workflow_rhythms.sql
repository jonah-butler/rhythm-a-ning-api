-- $2 must be an empty array (not NULL) when no blocks are kept
DELETE FROM workflow_rhythms
WHERE workflow_id = $1
  AND rhythm_id <> ALL($2::uuid[]);
