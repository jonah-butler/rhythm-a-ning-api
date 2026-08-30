INSERT INTO rhythm_tags (rhythm_id, tag_id)
SELECT $1, tag_id
FROM tags
WHERE name = ANY($2::text[])
ON CONFLICT DO NOTHING;
