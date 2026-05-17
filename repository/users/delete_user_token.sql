WITH deleted_count AS (
    DELETE FROM user_tokens 
    WHERE user_id = $1
      AND type = $2
      AND token = $3
    RETURNING *
)
SELECT EXISTS (SELECT count(*) FROM deleted_count) as deleted;