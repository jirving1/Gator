-- name: MarkFeedFetched :exec

UPDATE feeds_fetched
SET updated_at = CURRENT_TIMESTAMP, last_fetched_at =  CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2;

