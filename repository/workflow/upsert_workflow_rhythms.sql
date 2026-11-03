INSERT INTO workflow_rhythms (workflow_id, rhythm_id, section_id, measures, position)
SELECT $1::uuid, t.rhythm_id, t.section_id, t.measures, t.position
FROM UNNEST(
  $2::uuid[],     -- rhythm ids
  $3::uuid[],     -- section ids
  $4::smallint[], -- measures
  $5::smallint[]  -- positions
) AS t(rhythm_id, section_id, measures, position)
ON CONFLICT (workflow_id, rhythm_id) DO UPDATE
SET
  section_id = EXCLUDED.section_id,
  measures = EXCLUDED.measures,
  position = EXCLUDED.position;
