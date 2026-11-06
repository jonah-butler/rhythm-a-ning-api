SELECT
  r.rhythm_id,
  r.bpm,
  r.beats,
  st.name          AS subdivision,
  r.state,
  r.is_poly,
  r.poly_beats,
  pst.name         AS poly_subdivision,
  r.poly_state,
  r.user_id,
  l.name           AS level,
  r.name,
  r.description,
  r.sounds,
  r.poly_sounds,
  r.created_at,
  r.updated_at,
  r.subdivisions,
  r.poly_subdivisions,
  wr.section_id,
  wr.measures,
  wr.position
FROM workflow_rhythms wr
JOIN rhythms r ON r.rhythm_id = wr.rhythm_id
JOIN subdivision_types st ON st.subdivision_id = r.subdivision
LEFT JOIN subdivision_types pst ON pst.subdivision_id = r.poly_subdivision
LEFT JOIN levels l ON l.level_id = r.level_id
WHERE wr.workflow_id = $1
ORDER BY wr.position;
