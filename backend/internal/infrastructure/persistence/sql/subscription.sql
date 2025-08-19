-- name: getSubscriptionByIdForUser :many
SELECT s.id                    AS "subscriptions.id",
       s.owner_type            AS "subscriptions.owner_type",
       s.owner_family_id       AS "subscriptions.owner_family_id",
       s.owner_user_id         AS "subscriptions.owner_user_id",
       s.friendly_name         AS "subscriptions.friendly_name",
       s.free_trial_start_date AS "subscriptions.free_trial_start_date",
       s.free_trial_end_date   AS "subscriptions.free_trial_end_date",
       s.provider_id           AS "subscriptions.provider_id",
       s.plan_id               AS "subscriptions.plan_id",
       s.price_id              AS "subscriptions.price_id",
       s.family_id             AS "subscriptions.family_id",
       s.payer_type            AS "subscriptions.payer_type",
       s.payer_member_id       AS "subscriptions.payer_member_id",
       s.start_date            AS "subscriptions.start_date",
       s.end_date              AS "subscriptions.end_date",
       s.recurrency            AS "subscriptions.recurrency",
       s.custom_recurrency     AS "subscriptions.custom_recurrency",
       s.custom_price_currency AS "subscriptions.custom_price_currency",
       s.custom_price_amount   AS "subscriptions.custom_price_amount",
       s.created_at            AS "subscriptions.created_at",
       s.updated_at            AS "subscriptions.updated_at",
       s.etag                  AS "subscriptions.etag",
       su.family_member_id     AS "subscription_service_users.family_member_id"
FROM (SELECT s.*
      FROM public.subscriptions s
               LEFT JOIN subscription_service_users su ON su.subscription_id = s.id
               LEFT JOIN families f ON f.id = s.family_id
               LEFT JOIN family_members fm ON fm.id = su.family_member_id
      WHERE (s.owner_type = 'family' AND fm.user_id = $1)
         OR (s.owner_type = 'personal' AND s.owner_user_id = $1)) s
         LEFT JOIN subscription_service_users su ON su.subscription_id = s.id
WHERE s.id = $2;

-- name: getSubscriptionById :many
SELECT s.id                    AS "subscriptions.id",
       s.owner_type            AS "subscriptions.owner_type",
       s.owner_family_id       AS "subscriptions.owner_family_id",
       s.owner_user_id         AS "subscriptions.owner_user_id",
       s.friendly_name         AS "subscriptions.friendly_name",
       s.free_trial_start_date AS "subscriptions.free_trial_start_date",
       s.free_trial_end_date   AS "subscriptions.free_trial_end_date",
       s.provider_id           AS "subscriptions.provider_id",
       s.plan_id               AS "subscriptions.plan_id",
       s.price_id              AS "subscriptions.price_id",
       s.family_id             AS "subscriptions.family_id",
       s.payer_type            AS "subscriptions.payer_type",
       s.payer_member_id       AS "subscriptions.payer_member_id",
       s.start_date            AS "subscriptions.start_date",
       s.end_date              AS "subscriptions.end_date",
       s.recurrency            AS "subscriptions.recurrency",
       s.custom_recurrency     AS "subscriptions.custom_recurrency",
       s.custom_price_currency AS "subscriptions.custom_price_currency",
       s.custom_price_amount   AS "subscriptions.custom_price_amount",
       s.created_at            AS "subscriptions.created_at",
       s.updated_at            AS "subscriptions.updated_at",
       s.etag                  AS "subscriptions.etag",
       su.family_member_id     AS "subscription_service_users.family_member_id"
FROM public.subscriptions s
         LEFT JOIN subscription_service_users su ON su.subscription_id = s.id
WHERE s.id = $1;

-- name: getSubscriptionsForUser :many
WITH providers AS (SELECT p.id, p.name
                   FROM public.providers p
                   WHERE p.owner_type = 'system'
                      OR (p.owner_type = 'personal' AND p.owner_user_id = $1)
                      OR (p.owner_type = 'family' AND EXISTS (SELECT 1
                                                              FROM public.family_members fm
                                                              WHERE fm.family_id = p.owner_family_id
                                                                AND fm.user_id = $1))),
     matches AS (SELECT s.id
                 FROM public.subscriptions s
                 WHERE (
                     s.owner_type = 'system'
                         OR (s.owner_type = 'personal' AND s.owner_user_id = $1)
                         OR (s.owner_type = 'family' AND EXISTS (SELECT 1
                                                                 FROM public.family_members fm
                                                                 WHERE fm.family_id = s.owner_family_id
                                                                   AND fm.user_id = $1))
                     )
                   AND (
                     -- Match everything when the search term is NULL or empty
                     NULLIF(BTRIM($2), '') IS NULL
                         OR s.friendly_name ILIKE '%' || $2 || '%'
                         OR EXISTS (SELECT 1
                                    FROM providers p
                                    WHERE p.name ILIKE '%' || $2 || '%' AND p.id = s.provider_id)
                     )
                 ORDER BY $5
                 ),
     counted AS (SELECT m.id, COUNT(*) OVER () AS total_count
                 FROM matches m),
     paged AS (SELECT c.id, c.total_count
               FROM counted c
               ORDER BY c.id
               LIMIT $3 OFFSET $4)
SELECT s.id                    AS "subscriptions.id",
       s.owner_type            AS "subscriptions.owner_type",
       s.owner_family_id       AS "subscriptions.owner_family_id",
       s.owner_user_id         AS "subscriptions.owner_user_id",
       s.friendly_name         AS "subscriptions.friendly_name",
       s.free_trial_start_date AS "subscriptions.free_trial_start_date",
       s.free_trial_end_date   AS "subscriptions.free_trial_end_date",
       s.provider_id           AS "subscriptions.provider_id",
       s.plan_id               AS "subscriptions.plan_id",
       s.price_id              AS "subscriptions.price_id",
       s.family_id             AS "subscriptions.family_id",
       s.payer_type            AS "subscriptions.payer_type",
       s.payer_member_id       AS "subscriptions.payer_member_id",
       s.start_date            AS "subscriptions.start_date",
       s.end_date              AS "subscriptions.end_date",
       s.recurrency            AS "subscriptions.recurrency",
       s.custom_recurrency     AS "subscriptions.custom_recurrency",
       s.custom_price_currency AS "subscriptions.custom_price_currency",
       s.custom_price_amount   AS "subscriptions.custom_price_amount",
       s.created_at            AS "subscriptions.created_at",
       s.updated_at            AS "subscriptions.updated_at",
       s.etag                  AS "subscriptions.etag",
       su.family_member_id     AS "subscription_service_users.family_member_id",
       p.total_count           AS "total_count"
FROM paged p
         JOIN public.subscriptions s ON s.id = p.id
         LEFT JOIN subscription_service_users su ON su.subscription_id = s.id;

-- name: getSubscriptions :many
SELECT s.id                    AS "subscriptions.id",
       s.owner_type            AS "subscriptions.owner_type",
       s.owner_family_id       AS "subscriptions.owner_family_id",
       s.owner_user_id         AS "subscriptions.owner_user_id",
       s.friendly_name         AS "subscriptions.friendly_name",
       s.free_trial_start_date AS "subscriptions.free_trial_start_date",
       s.free_trial_end_date   AS "subscriptions.free_trial_end_date",
       s.provider_id           AS "subscriptions.provider_id",
       s.plan_id               AS "subscriptions.plan_id",
       s.price_id              AS "subscriptions.price_id",
       s.family_id             AS "subscriptions.family_id",
       s.payer_type            AS "subscriptions.payer_type",
       s.payer_member_id       AS "subscriptions.payer_member_id",
       s.start_date            AS "subscriptions.start_date",
       s.end_date              AS "subscriptions.end_date",
       s.recurrency            AS "subscriptions.recurrency",
       s.custom_recurrency     AS "subscriptions.custom_recurrency",
       s.custom_price_currency AS "subscriptions.custom_price_currency",
       s.custom_price_amount   AS "subscriptions.custom_price_amount",
       s.created_at            AS "subscriptions.created_at",
       s.updated_at            AS "subscriptions.updated_at",
       s.etag                  AS "subscriptions.etag",
       su.family_member_id     AS "subscription_service_users.family_member_id",
       s.total_count           AS "total_count"
FROM (SELECT s.*,
             COUNT(*) OVER () AS total_count
      FROM public.subscriptions s
      ORDER BY s.Id
      LIMIT $1 OFFSET $2) s
         LEFT JOIN subscription_service_users su ON su.subscription_id = s.id;

-- name: DeleteSubscription :exec
DELETE
FROM public.subscriptions
WHERE id = $1;

-- name: IsSubscriptionExists :one
SELECT COUNT(*)
FROM public.subscriptions s
WHERE s.id = ANY ($1::uuid[]);

-- name: CreateSubscription :exec
INSERT INTO public.subscriptions (id,
                                  friendly_name,
                                  free_trial_start_date,
                                  free_trial_end_date,
                                  provider_id,
                                  plan_id,
                                  price_id,
                                  custom_price_currency,
                                  custom_price_amount,
                                  owner_type,
                                  owner_family_id,
                                  owner_user_id,
                                  payer_type,
                                  family_id,
                                  payer_member_id,
                                  start_date,
                                  end_date,
                                  recurrency,
                                  custom_recurrency,
                                  created_at,
                                  updated_at,
                                  etag)
VALUES ($1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10, $11,
        $12, $13, $14, $15, $16, $17,
        $18, $19, $20, $21, $22);

-- name: CreateSubscriptionServiceUser :exec
INSERT INTO public.subscription_service_users (subscription_id, family_member_id)
VALUES ($1, $2);

-- name: CreateSubscriptions :copyfrom
INSERT INTO public.subscriptions (id,
                                  friendly_name,
                                  free_trial_start_date,
                                  free_trial_end_date,
                                  provider_id,
                                  plan_id,
                                  price_id,
                                  custom_price_currency,
                                  custom_price_amount,
                                  owner_type,
                                  owner_family_id,
                                  owner_user_id,
                                  payer_type,
                                  family_id,
                                  payer_member_id,
                                  start_date,
                                  end_date,
                                  recurrency,
                                  custom_recurrency,
                                  created_at,
                                  updated_at,
                                  etag)
VALUES ($1, $2, $3, $4, $5, $6, $7,
        $8, $9, $10, $11, $12,
        $13, $14, $15, $16, $17, $18,
        $19, $20, $21, $22);

-- name: CreateSubscriptionServiceUsers :copyfrom
INSERT INTO public.subscription_service_users (subscription_id, family_member_id)
VALUES ($1, $2);

-- name: UpdateSubscription :exec
UPDATE public.subscriptions
SET friendly_name         = $2,
    free_trial_start_date = $3,
    free_trial_end_date   = $4,
    provider_id           = $5,
    plan_id               = $6,
    price_id              = $7,
    custom_price_currency = $8,
    custom_price_amount   = $9,
    owner_type            = $10,
    owner_family_id       = $11,
    owner_user_id         = $12,
    payer_type            = $13,
    family_id             = $14,
    payer_member_id       = $15,
    start_date            = $16,
    end_date              = $17,
    recurrency            = $18,
    custom_recurrency     = $19,
    updated_at            = $20,
    etag                  = $21
WHERE id = $1;

-- name: DeleteSubscriptionServiceUser :exec
DELETE
FROM public.subscription_service_users
WHERE subscription_id = $1
  AND family_member_id = $2;