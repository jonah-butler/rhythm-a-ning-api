WITH inserted AS (
  INSERT INTO users (username, email, password)
  VALUES ($1, $2, $3)
  ON CONFLICT DO NOTHING
  RETURNING user_id, email
)

SELECT
  (SELECT user_id FROM inserted) AS user_id,
  (SELECT email FROM inserted) AS email,
  EXISTS(SELECT 1 FROM users WHERE username = $1) AS username_taken,
  EXISTS(SELECT 1 FROM users WHERE email = $2) as email_taken;