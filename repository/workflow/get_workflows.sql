-- totals mirror the client's totalMeasures, averageTempo and totalTime (minutes)
SELECT
  w.workflow_id,
  w.user_id,
  w.name,
  w.description,
  w.created_at,
  w.updated_at,
  COUNT(wr.rhythm_id)                                              AS block_count,
  COALESCE(SUM(wr.measures), 0)                                    AS total_measures,
  COALESCE(FLOOR(AVG(r.bpm)), 0)::int                              AS average_tempo,
  COALESCE(SUM(r.beats::numeric / r.bpm * wr.measures), 0)::float8 AS total_minutes,
  COUNT(*) OVER()                                                  AS total
FROM workflows w
LEFT JOIN workflow_rhythms wr ON wr.workflow_id = w.workflow_id
LEFT JOIN rhythms r ON r.rhythm_id = wr.rhythm_id
WHERE w.user_id = $1
GROUP BY w.workflow_id
ORDER BY w.updated_at DESC
LIMIT 10 OFFSET 10 * $2;
