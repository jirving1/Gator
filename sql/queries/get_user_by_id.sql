-- name: GetUserByID :one

SELECT name FROM users 
WHERE id = $1;