WITH inserted AS (
  INSERT INTO rhythms (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    poly_sounds,
    user_id,
    name,
    description,
    level_id,
    sounds
  )
  VALUES (
    $1, -- bpm
    $2, -- beats
    (SELECT subdivision_id FROM subdivision_types WHERE name = $3), -- subdivision
    $4, -- state
    $5, -- is poly
    $6, --poly beats
    (SELECT subdivision_id FROM subdivision_types WHERE name = $7), -- poly subdivision
    $8, -- poly state
    $9, -- poly sounds
    $10, -- user id
    $11, -- name
    $12, -- description
    (SELECT level_id FROM levels WHERE name = $13), -- level
    $14
  )
  RETURNING *
)
SELECT
  inserted.rhythm_id AS id,
  inserted.bpm,
  inserted.beats,
  st.name AS subdivision, -- shadows subdivision with string value used in client
  inserted.state,
  inserted.is_poly,
  inserted.poly_beats,
  pst.name AS poly_subdivision, -- shadows poly subdivision with string || NULL value used in client
  inserted.poly_state,
  inserted.poly_sounds,
  inserted.user_id,
  inserted.created_at,
  inserted.updated_at,
  inserted.name,
  inserted.description,
  l.name as level,
  inserted.sounds
FROM inserted
JOIN subdivision_types st ON st.subdivision_id = inserted.subdivision
LEFT JOIN subdivision_types pst ON pst.subdivision_id = inserted.poly_subdivision
LEFT JOIN levels l ON l.level_id = inserted.level_id;
       