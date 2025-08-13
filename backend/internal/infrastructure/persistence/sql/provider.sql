-- name: getProviderById :many
SELECT p.id               AS "providers.id",
       p.owner_type       AS "providers.owner_type",
       p.owner_family_id  AS "providers.owner_family_id",
       p.owner_user_id    AS "providers.owner_user_id",
       p.name             AS "providers.name",
       p.key              AS "providers.key",
       p.description      AS "providers.description",
       p.icon_url         AS "providers.icon_url",
       p.url              AS "providers.url",
       p.pricing_page_url AS "providers.pricing_page_url",
       p.created_at       AS "providers.created_at",
       p.updated_at       AS "providers.updated_at",
       p.etag             AS "providers.etag",
       ppl.id             AS "provider_plans.id",
       ppl.name           AS "provider_plans.name",
       ppl.description    AS "provider_plans.description",
       ppl.provider_id    AS "provider_plans.provider_id",
       ppl.created_at     AS "provider_plans.created_at",
       ppl.updated_at     AS "provider_plans.updated_at",
       ppl.etag           AS "provider_plans.etag",
       ppr.id             AS "provider_prices.id",
       ppr.start_date     AS "provider_prices.start_date",
       ppr.end_date       AS "provider_prices.end_date",
       ppr.currency       AS "provider_prices.currency",
       ppr.amount         AS "provider_prices.amount",
       ppr.plan_id        AS "provider_prices.plan_id",
       ppr.created_at     AS "provider_prices.created_at",
       ppr.updated_at     AS "provider_prices.updated_at",
       ppr.etag           AS "provider_prices.etag",
       pl.label_id        AS "provider_labels.label_id",
       pl.provider_id     AS "provider_labels.provider_id"
FROM public.providers p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
         LEFT JOIN public.provider_labels pl ON pl.provider_id = p.id
WHERE p.id = $1;

-- name: getProviderByIdForUser :many
SELECT p.id               AS "providers.id",
       p.owner_type       AS "providers.owner_type",
       p.owner_family_id  AS "providers.owner_family_id",
       p.owner_user_id    AS "providers.owner_user_id",
       p.name             AS "providers.name",
       p.key              AS "providers.key",
       p.description      AS "providers.description",
       p.icon_url         AS "providers.icon_url",
       p.url              AS "providers.url",
       p.pricing_page_url AS "providers.pricing_page_url",
       p.created_at       AS "providers.created_at",
       p.updated_at       AS "providers.updated_at",
       p.etag             AS "providers.etag",
       ppl.id             AS "provider_plans.id",
       ppl.name           AS "provider_plans.name",
       ppl.description    AS "provider_plans.description",
       ppl.provider_id    AS "provider_plans.provider_id",
       ppl.created_at     AS "provider_plans.created_at",
       ppl.updated_at     AS "provider_plans.updated_at",
       ppl.etag           AS "provider_plans.etag",
       ppr.id             AS "provider_prices.id",
       ppr.start_date     AS "provider_prices.start_date",
       ppr.end_date       AS "provider_prices.end_date",
       ppr.currency       AS "provider_prices.currency",
       ppr.amount         AS "provider_prices.amount",
       ppr.plan_id        AS "provider_prices.plan_id",
       ppr.created_at     AS "provider_prices.created_at",
       ppr.updated_at     AS "provider_prices.updated_at",
       ppr.etag           AS "provider_prices.etag",
       pl.label_id        AS "provider_labels.label_id",
       pl.provider_id     AS "provider_labels.provider_id"
FROM public.providers p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
         LEFT JOIN public.provider_labels pl ON pl.provider_id = p.id
         LEFT JOIN public.families f ON f.id = p.owner_family_id
         LEFT JOIN public.family_members fm ON fm.family_id = f.id
WHERE p.id = $1
  AND (
    p.owner_type = 'system' OR
    (p.owner_type = 'personal' AND p.owner_user_id = $2) OR
    (p.owner_type = 'family' AND fm.user_id = $2)
    );

-- name: getSystemProviders :many
SELECT p.id               AS "providers.id",
       p.owner_type       AS "providers.owner_type",
       p.owner_family_id  AS "providers.owner_family_id",
       p.owner_user_id    AS "providers.owner_user_id",
       p.name             AS "providers.name",
       p.key              AS "providers.key",
       p.description      AS "providers.description",
       p.icon_url         AS "providers.icon_url",
       p.url              AS "providers.url",
       p.pricing_page_url AS "providers.pricing_page_url",
       p.created_at       AS "providers.created_at",
       p.updated_at       AS "providers.updated_at",
       p.etag             AS "providers.etag",
       ppl.id             AS "provider_plans.id",
       ppl.name           AS "provider_plans.name",
       ppl.description    AS "provider_plans.description",
       ppl.provider_id    AS "provider_plans.provider_id",
       ppl.created_at     AS "provider_plans.created_at",
       ppl.updated_at     AS "provider_plans.updated_at",
       ppl.etag           AS "provider_plans.etag",
       ppr.id             AS "provider_prices.id",
       ppr.start_date     AS "provider_prices.start_date",
       ppr.end_date       AS "provider_prices.end_date",
       ppr.currency       AS "provider_prices.currency",
       ppr.amount         AS "provider_prices.amount",
       ppr.plan_id        AS "provider_prices.plan_id",
       ppr.created_at     AS "provider_prices.created_at",
       ppr.updated_at     AS "provider_prices.updated_at",
       ppr.etag           AS "provider_prices.etag",
       pl.label_id        AS "provider_labels.label_id",
       pl.provider_id     AS "provider_labels.provider_id",
       COUNT(*) OVER ()   AS total_count
FROM public.providers p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
         LEFT JOIN public.provider_labels pl ON pl.provider_id = p.id
WHERE p.owner_type = 'system'
  AND p.owner_user_id IS NULL
  AND p.owner_family_id IS NULL;

-- name: getProviders :many
SELECT p.id               AS "providers.id",
       p.owner_type       AS "providers.owner_type",
       p.owner_family_id  AS "providers.owner_family_id",
       p.owner_user_id    AS "providers.owner_user_id",
       p.name             AS "providers.name",
       p.key              AS "providers.key",
       p.description      AS "providers.description",
       p.icon_url         AS "providers.icon_url",
       p.url              AS "providers.url",
       p.pricing_page_url AS "providers.pricing_page_url",
       p.created_at       AS "providers.created_at",
       p.updated_at       AS "providers.updated_at",
       p.etag             AS "providers.etag",
       ppl.id             AS "provider_plans.id",
       ppl.name           AS "provider_plans.name",
       ppl.description    AS "provider_plans.description",
       ppl.provider_id    AS "provider_plans.provider_id",
       ppl.created_at     AS "provider_plans.created_at",
       ppl.updated_at     AS "provider_plans.updated_at",
       ppl.etag           AS "provider_plans.etag",
       ppr.id             AS "provider_prices.id",
       ppr.start_date     AS "provider_prices.start_date",
       ppr.end_date       AS "provider_prices.end_date",
       ppr.currency       AS "provider_prices.currency",
       ppr.amount         AS "provider_prices.amount",
       ppr.plan_id        AS "provider_prices.plan_id",
       ppr.created_at     AS "provider_prices.created_at",
       ppr.updated_at     AS "provider_prices.updated_at",
       ppr.etag           AS "provider_prices.etag",
       pl.label_id        AS "provider_labels.label_id",
       pl.provider_id     AS "provider_labels.provider_id",
       p.total_count      AS "total_count"
FROM (SELECT *,
             COUNT(*) OVER () AS total_count
      FROM public.providers
      ORDER BY id
      LIMIT $1 OFFSET $2) p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
         LEFT JOIN public.provider_labels pl ON pl.provider_id = p.id;

-- name: getProvidersForUser :many
SELECT p.id               AS "providers.id",
       p.owner_type       AS "providers.owner_type",
       p.owner_family_id  AS "providers.owner_family_id",
       p.owner_user_id    AS "providers.owner_user_id",
       p.name             AS "providers.name",
       p.key              AS "providers.key",
       p.description      AS "providers.description",
       p.icon_url         AS "providers.icon_url",
       p.url              AS "providers.url",
       p.pricing_page_url AS "providers.pricing_page_url",
       p.created_at       AS "providers.created_at",
       p.updated_at       AS "providers.updated_at",
       p.etag             AS "providers.etag",
       ppl.id             AS "provider_plans.id",
       ppl.name           AS "provider_plans.name",
       ppl.description    AS "provider_plans.description",
       ppl.provider_id    AS "provider_plans.provider_id",
       ppl.created_at     AS "provider_plans.created_at",
       ppl.updated_at     AS "provider_plans.updated_at",
       ppl.etag           AS "provider_plans.etag",
       ppr.id             AS "provider_prices.id",
       ppr.start_date     AS "provider_prices.start_date",
       ppr.end_date       AS "provider_prices.end_date",
       ppr.currency       AS "provider_prices.currency",
       ppr.amount         AS "provider_prices.amount",
       ppr.plan_id        AS "provider_prices.plan_id",
       ppr.created_at     AS "provider_prices.created_at",
       ppr.updated_at     AS "provider_prices.updated_at",
       ppr.etag           AS "provider_prices.etag",
       pl.label_id        AS "provider_labels.label_id",
       pl.provider_id     AS "provider_labels.provider_id",
       p.total_count      AS "total_count"
FROM (SELECT p.*,
             COUNT(*) OVER () AS total_count
      FROM public.providers p
               LEFT JOIN public.families f ON f.id = p.owner_family_id
               LEFT JOIN public.family_members fm ON fm.family_id = f.id
      WHERE p.owner_type = 'system'
         OR (p.owner_type = 'personal' AND p.owner_user_id = $1)
         OR (p.owner_type = 'family' AND fm.user_id = $1)
      ORDER BY p.id
      LIMIT $2 OFFSET $3) p
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = p.id
         LEFT JOIN public.provider_prices ppr ON ppl.provider_id = p.id
         LEFT JOIN public.provider_labels pl ON pl.provider_id = p.id;

-- test
WITH matches AS (
    SELECT p.id
    FROM public.providers p
    WHERE
        (
            p.owner_type = 'system'
                OR (p.owner_type = 'personal' AND p.owner_user_id = 'kp_4de08ac140e54e5081b192b48c1b660e')
                OR (p.owner_type = 'family' AND EXISTS (
                SELECT 1
                FROM public.family_members fm
                WHERE fm.family_id = p.owner_family_id
                  AND fm.user_id = 'kp_4de08ac140e54e5081b192b48c1b660e'
            ))
            )
      AND (
        -- Match everything when the search term is NULL or empty
        NULLIF(BTRIM('red'), '') IS NULL
            OR p.name ILIKE '%' || 'red' || '%'
            OR EXISTS (
            SELECT 1
            FROM public.provider_labels pl
                     JOIN public.labels l ON l.id = pl.label_id
            WHERE pl.provider_id = p.id
              AND l.name ILIKE '%' || 'red' || '%'
        )
        )
),
     counted AS (
         SELECT m.id, COUNT(*) OVER () AS total_count
         FROM matches m
     ),
     paged AS (
         SELECT c.id, c.total_count
         FROM counted c
         ORDER BY c.id
         LIMIT 10 OFFSET 0
     )
SELECT pr.id               AS "providers.id",
       pr.owner_type       AS "providers.owner_type",
       pr.owner_family_id  AS "providers.owner_family_id",
       pr.owner_user_id    AS "providers.owner_user_id",
       pr.name             AS "providers.name",
       pr.key              AS "providers.key",
       pr.description      AS "providers.description",
       pr.icon_url         AS "providers.icon_url",
       pr.url              AS "providers.url",
       pr.pricing_page_url AS "providers.pricing_page_url",
       pr.created_at       AS "providers.created_at",
       pr.updated_at       AS "providers.updated_at",
       pr.etag             AS "providers.etag",
       ppl.id              AS "provider_plans.id",
       ppl.name            AS "provider_plans.name",
       ppl.description     AS "provider_plans.description",
       ppl.provider_id     AS "provider_plans.provider_id",
       ppl.created_at      AS "provider_plans.created_at",
       ppl.updated_at      AS "provider_plans.updated_at",
       ppl.etag            AS "provider_plans.etag",
       ppr.id              AS "provider_prices.id",
       ppr.start_date      AS "provider_prices.start_date",
       ppr.end_date        AS "provider_prices.end_date",
       ppr.currency        AS "provider_prices.currency",
       ppr.amount          AS "provider_prices.amount",
       ppr.plan_id         AS "provider_prices.plan_id",
       ppr.created_at      AS "provider_prices.created_at",
       ppr.updated_at      AS "provider_prices.updated_at",
       ppr.etag            AS "provider_prices.etag",
       pl.label_id         AS "provider_labels.label_id",
       pl.provider_id      AS "provider_labels.provider_id",
       p.total_count       AS "total_count"
FROM paged p
         JOIN public.providers pr ON pr.id = p.id
         LEFT JOIN public.provider_plans ppl ON ppl.provider_id = pr.id
         LEFT JOIN public.provider_prices ppr ON ppr.plan_id = ppl.id
         LEFT JOIN public.provider_labels pl ON pl.provider_id = pr.id;


-- name: getProvidersForUserWithSearch :many
WITH matches AS (
    SELECT p.id
    FROM public.providers p
    WHERE
        (
            p.owner_type = 'system'
         OR (p.owner_type = 'personal' AND p.owner_user_id = $1)
         OR (p.owner_type = 'family' AND EXISTS (
                SELECT 1
                FROM public.family_members fm
                WHERE fm.family_id = p.owner_family_id
                  AND fm.user_id = $1
            ))
        )
        AND (
            p.name ILIKE $2
         OR EXISTS (
                SELECT 1
                FROM public.provider_labels pl
                JOIN public.labels l ON l.id = pl.label_id
                WHERE pl.provider_id = p.id
                  AND l.name ILIKE $2
            )
        )
),
counted AS (
    SELECT m.id, COUNT(*) OVER () AS total_count
    FROM matches m
),
paged AS (
    SELECT c.id, c.total_count
    FROM counted c
    ORDER BY c.id
    LIMIT $3 OFFSET $4
)
SELECT pr.id               AS "providers.id",
       pr.owner_type       AS "providers.owner_type",
       pr.owner_family_id  AS "providers.owner_family_id",
       pr.owner_user_id    AS "providers.owner_user_id",
       pr.name             AS "providers.name",
       pr.key              AS "providers.key",
       pr.description      AS "providers.description",
       pr.icon_url         AS "providers.icon_url",
       pr.url              AS "providers.url",
       pr.pricing_page_url AS "providers.pricing_page_url",
       pr.created_at       AS "providers.created_at",
       pr.updated_at       AS "providers.updated_at",
       pr.etag             AS "providers.etag",
       ppl.id              AS "provider_plans.id",
       ppl.name            AS "provider_plans.name",
       ppl.description     AS "provider_plans.description",
       ppl.provider_id     AS "provider_plans.provider_id",
       ppl.created_at      AS "provider_plans.created_at",
       ppl.updated_at      AS "provider_plans.updated_at",
       ppl.etag            AS "provider_plans.etag",
       ppr.id              AS "provider_prices.id",
       ppr.start_date      AS "provider_prices.start_date",
       ppr.end_date        AS "provider_prices.end_date",
       ppr.currency        AS "provider_prices.currency",
       ppr.amount          AS "provider_prices.amount",
       ppr.plan_id         AS "provider_prices.plan_id",
       ppr.created_at      AS "provider_prices.created_at",
       ppr.updated_at      AS "provider_prices.updated_at",
       ppr.etag            AS "provider_prices.etag",
       pl.label_id         AS "provider_labels.label_id",
       pl.provider_id      AS "provider_labels.provider_id",
       p.total_count       AS "total_count"
FROM paged p
JOIN public.providers pr ON pr.id = p.id
LEFT JOIN public.provider_plans ppl ON ppl.provider_id = pr.id
LEFT JOIN public.provider_prices ppr ON ppr.plan_id = ppl.id
LEFT JOIN public.provider_labels pl ON pl.provider_id = pr.id;


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

-- name: CreateProviderLabel :exec
INSERT INTO public.provider_labels (label_id, provider_id)
VALUES ($1, $2);

-- name: CreateProviderLabels :copyfrom
INSERT INTO public.provider_labels (label_id, provider_id)
VALUES ($1, $2);

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

-- name: DeleteProviderPlan :exec
DELETE
FROM public.provider_plans
WHERE id = $1;

-- name: DeleteProviderPrice :exec
DELETE
FROM public.provider_prices
WHERE id = $1;

-- name: DeleteProviderLabel :exec
DELETE
FROM public.provider_labels
WHERE provider_id = $1
  AND label_id = $2;

-- name: IsProviderExists :one
SELECT COUNT(*)
FROM public.providers p
WHERE p.id = ANY ($1::uuid[]);