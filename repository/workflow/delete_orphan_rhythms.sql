-- removes rhythms that were in the workflow before the update
-- and are no longer linked to any workflow
DELETE FROM rhythms r
WHERE r.rhythm_id = ANY($1::uuid[])
  AND r.user_id = $2
  AND NOT EXISTS (
    SELECT 1 FROM workflow_rhythms wr WHERE wr.rhythm_id = r.rhythm_id
  );
