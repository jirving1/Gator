-- name: GetFeed :one 

SELECT * FROM feeds 
WHERE name = $1;

-- name: GetFeeds :many

SELECT * FROM feeds;

-- name: GetFeedByID :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetFeedByURL :one

SELECT * FROM feeds 
WHERE url = $1;