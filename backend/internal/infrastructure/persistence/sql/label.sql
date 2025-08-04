-- name: GetLabelById :one
SELECT *
FROM public.labels l
WHERE l.id = $1;

-- name: GetLabels :many
SELECT sqlc.embed(l),
       COUNT(*) OVER () AS total_count
FROM labels l
WHERE (
          -- Personal owner condition
          ('personal' = ANY($1::varchar[]) AND l.owner_user_id = $2)
              OR
              -- System owner condition
          ('system'= ANY($1::varchar[]))
              OR
              -- Family owner condition
          ('family' = ANY($1::varchar[]) AND l.owner_family_id = ANY ($3::uuid[]))
          )
LIMIT $4 OFFSET $5;

-- name: GetSystemLabels :many
SELECT *
FROM public.labels l
WHERE l.owner_family_id IS NULL
  AND l.owner_user_id IS NULL
  AND l.owner_type = 'system';

-- name: IsLabelExists :one
SELECT COUNT(*)
FROM public.labels l
WHERE l.id = ANY ($1::uuid[]);

-- name: CreateLabel :exec
INSERT INTO public.labels (id, owner_type, owner_family_id, owner_user_id, name, key, color, created_at, updated_at,
                           etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: CreateLabels :copyfrom
INSERT INTO public.labels (id, owner_type, owner_family_id, owner_user_id, name, key, color, created_at, updated_at,
                           etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdateLabel :exec
UPDATE public.labels
SET owner_type      = $2,
    owner_family_id = $3,
    owner_user_id   = $4,
    name            = $5,
    key             = $6,
    color           = $7,
    updated_at      = $8,
    etag            = $9
WHERE id = $1;

-- name: DeleteLabel :exec
DELETE
FROM public.labels l
WHERE l.id = $1;

-- name: DeleteLabelsForUser :exec
DELETE
FROM public.labels l
WHERE l.owner_type = 'user'
  AND l.owner_user_id = $1;

-- name: DeleteLabelsForFamily :exec
DELETE
FROM public.labels l
WHERE l.owner_type = 'family'
  AND l.owner_family_id = $1;