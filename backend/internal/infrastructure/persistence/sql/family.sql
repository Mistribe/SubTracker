-- name: getFamilyById :many
SELECT f.id          AS "families.id",
       f.name        AS "families.name",
       f.owner_id    AS "families.owner_id",
       f.created_at  AS "families.created_at",
       f.updated_at  AS "families.updated_at",
       f.etag        AS "families.etag",
       fm.id         AS "family_members.id",
       fm.name       AS "family_members.name",
       fm.family_id  AS "family_members.family_id",
       fm.user_id    AS "family_members.user_id",
       fm.type       AS "family_members.type",
       fm.created_at AS "family_members.created_at",
       fm.updated_at AS "family_members.updated_at",
       fm.etag       AS "family_members.etag"
FROM public.families f
         LEFT JOIN public.family_members fm ON f.id = fm.family_id
WHERE f.id = $1;

-- name: getFamiliesForUser :many
SELECT f.id          AS "families.id",
       f.name        AS "families.name",
       f.owner_id    AS "families.owner_id",
       f.created_at  AS "families.created_at",
       f.updated_at  AS "families.updated_at",
       f.etag        AS "families.etag",
       fm.id         AS "family_members.id",
       fm.name       AS "family_members.name",
       fm.family_id  AS "family_members.family_id",
       fm.user_id    AS "family_members.user_id",
       fm.type       AS "family_members.type",
       fm.created_at AS "family_members.created_at",
       fm.updated_at AS "family_members.updated_at",
       fm.etag       AS "family_members.etag",
       f.total_count AS "total_count"
FROM (SELECT f.*,
             COUNT(*) OVER () AS total_count
      FROM public.families f
               LEFT JOIN public.family_members fm ON f.id = fm.family_id
      WHERE fm.user_id = $1
      ORDER BY f.Id
      LIMIT $2 OFFSET $3) f
         LEFT JOIN public.family_members fm ON f.id = fm.family_id;

-- name: IsMemberOfFamily :one
SELECT COUNT(*)
FROM public.family_members fm
WHERE fm.family_id = $1
  AND fm.user_id = $2;

-- name: IsMemberExists :one
SELECT COUNT(*)
FROM public.family_members fm
WHERE fm.family_id = $1
  AND fm.id = ANY ($2::uuid[]);

-- name: IsFamilyExists :one
SELECT COUNT(*)
FROM public.families f
WHERE f.id = ANY ($1::uuid[]);

-- name: CreateFamily :exec
INSERT INTO public.families (id, name, owner_id, created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateFamilies :copyfrom
INSERT INTO public.families (id, name, owner_id, created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateFamilyMember :exec
INSERT INTO public.family_members (id, family_id, user_id, name, type, created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateFamilyMembers :copyfrom
INSERT INTO public.family_members (id, family_id, user_id, name, type, created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateFamily :exec
UPDATE public.families
SET name       = $2,
    owner_id   = $3,
    updated_at = $4,
    etag       = $5
WHERE id = $1;

-- name: UpdateFamilyMember :exec
UPDATE public.family_members
SET family_id  = $2,
    user_id    = $3,
    name       = $4,
    type       = $5,
    updated_at = $6,
    etag       = $7
WHERE id = $1;

-- name: DeleteFamily :exec
DELETE
FROM public.families f
WHERE f.id = $1;

-- name: DeleteFamilyMember :exec
DELETE
FROM public.family_members fm
WHERE fm.id = $1;