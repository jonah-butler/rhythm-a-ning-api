WITH deleted as (
  DELETE FROM rhythms
  WHERE rhythm_id = $1
  AND user_id = $2
  RETURNING *
)
SELECT EXISTS (SELECT 1 from deleted) as deleted;