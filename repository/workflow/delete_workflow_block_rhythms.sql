-- deleting the rhythms cascades to their workflow_rhythms rows
DELETE FROM rhythms r
USING workflow_rhythms wr
WHERE wr.rhythm_id = r.rhythm_id
  AND wr.workflow_id = $1
  AND r.user_id = $2;
