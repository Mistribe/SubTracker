-- name: GetProviderById :many
SELECT sqlc.embed(p),
       sqlc.embed(ppl),
       sqlc.embed(ppr)
FROM public.providers p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
WHERE p.id = $1;

-- name: GetSystemProviders :many
SELECT sqlc.embed(p),
       sqlc.embed(ppl),
       sqlc.embed(ppr)
FROM public.providers p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
WHERE p.owner_type = 'system'
  AND p.owner_user_id IS NULL
  AND p.owner_family_id IS NULL;

-- name: GetProviders :many
SELECT sqlc.embed(p),
       sqlc.embed(ppl),
       sqlc.embed(ppr),
       COUNT() OVER () AS total_count
FROM public.providers p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
LIMIT $1 OFFSET $2;


-- name: CreateProviderPrice :exec
INSERT INTO public.provider_prices (id, plan_id, currency,
                                    amount, start_date, end_date,
                                    created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: CreateProviderPrices :copyfrom
INSERT INTO public.provider_prices (id, plan_id, currency,
                                    amount, start_date, end_date,
                                    created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: CreateProviderPlan :exec
INSERT INTO public.provider_plans (id, provider_id, name,
                                   description, created_at, updated_at,
                                   etag)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateProviderPlans :copyfrom
INSERT INTO public.provider_plans (id, provider_id, name,
                                   description, created_at, updated_at,
                                   etag)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateProvider :exec
INSERT INTO public.providers (id, owner_type, owner_family_id,
                              owner_user_id, name, key,
                              description, icon_url, url,
                              pricing_page_url, created_at, updated_at,
                              etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: CreateProviders :copyfrom
INSERT INTO public.providers (id, owner_type, owner_family_id,
                              owner_user_id, name, key,
                              description, icon_url, url,
                              pricing_page_url, created_at, updated_at,
                              etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: UpdateProviderPrice :exec
UPDATE public.provider_prices
SET currency   = $2,
    amount     = $3,
    start_date = $4,
    end_date   = $5,
    updated_at = $6,
    etag       = $7,
    plan_id    = $8
WHERE id = $1;

-- name: UpdateProviderPlan :exec
UPDATE public.provider_plans
SET name        = $2,
    description = $3,
    updated_at  = $4,
    etag        = $5,
    provider_id = $6
WHERE id = $1;

-- name: UpdateProvider :exec
UPDATE public.providers
SET owner_type       = $2,
    owner_family_id  = $3,
    owner_user_id    = $4,
    name             = $5,
    key              = $6,
    description      = $7,
    icon_url         = $8,
    url              = $9,
    pricing_page_url = $10,
    updated_at       = $11
WHERE id = $1;

-- name: DeleteProvider :exec
DELETE
FROM public.providers
WHERE id = $1;