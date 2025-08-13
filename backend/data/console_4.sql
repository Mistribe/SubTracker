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
        NULLIF(BTRIM($1), '') IS NULL
        OR p.name ILIKE '%' || $1 || '%'
        OR EXISTS (
            SELECT 1
            FROM public.provider_labels pl
            JOIN public.labels l ON l.id = pl.label_id
            WHERE pl.provider_id = p.id
              AND l.name ILIKE '%' || $1 || '%'
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
