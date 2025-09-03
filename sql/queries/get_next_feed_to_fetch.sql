-- name: GetNextFeedToFetch :one
SELECT * FROM feed_follows
WHERE user_id = $1 
ORDER BY updated_at ASC NULLS FIRST;