-- name: ListPackSizes :many
SELECT size
FROM pack_sizes
ORDER BY size ASC;

-- name: DeletePackSizes :exec
DELETE FROM pack_sizes;

-- name: CreatePackSize :exec
INSERT INTO pack_sizes (size)
VALUES ($1);

-- name: DeletePackSize :exec
DELETE FROM pack_sizes
WHERE size = $1;
