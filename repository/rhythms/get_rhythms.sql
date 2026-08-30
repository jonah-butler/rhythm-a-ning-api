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
    l.name as level,
    r.name,
    r.description,
    r.sounds,
    r.poly_sounds,
    r.created_at,
    r.updated_at,
    COUNT(*) OVER() as total
FROM rhythms r
JOIN subdivision_types st ON st.subdivision_id = r.subdivision
LEFT JOIN subdivision_types pst ON pst.subdivision_id = r.poly_subdivision
LEFT JOIN levels l ON l.level_id = r.level_id
WHERE r.user_id = $1
ORDER BY r.created_at DESC
LIMIT 10 OFFSET 10 * $2;