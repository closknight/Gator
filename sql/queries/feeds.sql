-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name as feed_name, url, users.name as user_name from feeds
INNER JOIN users
ON users.id = feeds.user_id;

-- name: GetFeedByURL :one
SELECT * from feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds 
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
