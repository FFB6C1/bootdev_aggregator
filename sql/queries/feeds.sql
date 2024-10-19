-- name: AddFeed :exec
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * 
FROM feeds;

-- name: GetFeedIDByURL :one
SELECT id
FROM feeds
WHERE url = $1;

-- name: MarkFetchedByURL :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE url = $1;

-- name: GetNextFeed :one
SELECT feeds.url AS url
FROM feeds
INNER JOIN feed_follows
ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;