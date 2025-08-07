-- name: GetUserFamilyIds :many
SELECT fm.family_id FROM public.family_members fm
where fm.user_id = $1;

-- name: GetUserProfile :one
SELECT * from public.users u
where u.id = $1;

-- name: CreateUserProfile :exec
INSERT INTO public.users (id, currency)
VALUES ($1, $2)
ON CONFLICT (id)
    DO UPDATE SET currency = EXCLUDED.currency;
