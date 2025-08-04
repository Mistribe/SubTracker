-- name: GetSubscriptionById :many
SELECT sqlc.embed(s),
       sqlc.embed(su)
FROM public.subscriptions s
         INNER JOIN subscription_service_users su ON su.subscription_id = s.id
WHERE s.id = $1;

-- name: GetSubscriptions :many
SELECT sqlc.embed(s),
       sqlc.embed(su),
       COUNT(*) OVER () AS total_count
FROM public.subscriptions s
         INNER JOIN subscription_service_users su ON su.subscription_id = s.id
LIMIT $1 OFFSET $2;

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