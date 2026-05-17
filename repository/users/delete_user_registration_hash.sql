WITH deleted_count AS (
    DELETE FROM user_tokens 
    WHERE user_id = $1
      AND type = 1
      AND token = $2
    RETURNING *
)
SELECT EXISTS (SELECT count(*) FROM deleted_count) as did_delete;