WITH updated AS (
  UPDATE rhythms
  SET
    bpm              = $1,
    beats            = $2,
    subdivision      = (SELECT subdivision_id FROM subdivision_types WHERE name = $3),
    state            = $4,
    is_poly          = $5,
    poly_beats       = $6,
    poly_subdivision = (SELECT subdivision_id FROM subdivision_types WHERE name = $7),
    poly_state       = $8,
    level_id         = (SELECT level_id FROM levels WHERE name = $9),
    name             = $10,
    description      = $11,       
    updated_at       = NOW()
  WHERE rhythm_id = $12
  AND user_id     = $13
  RETURNING *
)
SELECT
  updated.rhythm_id,
  updated.bpm,
  updated.beats,
  st.name          AS subdivision,
  updated.state,
  updated.is_poly,
  updated.poly_beats,
  pst.name         AS poly_subdivision,
  updated.poly_state,
  updated.user_id,
  updated.name,
  updated.description,
  l.name as level,
  updated.created_at,
  updated.updated_at
FROM updated
JOIN subdivision_types st ON st.subdivision_id = updated.subdivision
LEFT JOIN subdivision_types pst ON pst.subdivision_id = updated.poly_subdivision
LEFT JOIN levels l ON l.level_id = updated.level_id;