-- name: GetPostsForUser :many
SELECT
  p.*,
  f.name AS feed_name
FROM posts p
JOIN feed_follows ff ON ff.feed_id = p.feed_id
JOIN feeds f ON f.id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC NULLS LAST, p.created_at DESC
LIMIT $2;

