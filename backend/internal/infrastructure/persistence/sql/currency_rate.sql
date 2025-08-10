-- name: GetCurrencyRateById :one
SELECT sqlc.embed(cr)
FROM public.currency_rates cr
WHERE cr.id = $1;

-- name: GetCurrencyRatesByDate :many
SELECT sqlc.embed(cr)
FROM public.currency_rates cr
WHERE cr.rate_date = $1
ORDER BY cr.from_currency, cr.to_currency;

-- name: GetLatestUpdateDate :one
SELECT COALESCE(MAX(cr.updated_at), 'epoch'::timestamp)::timestamp AS latest_update_date
FROM public.currency_rates cr;

-- name: CreateCurrencyRate :exec
INSERT INTO public.currency_rates (id, from_currency, to_currency, rate_date, exchange_rate,
                                   created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateCurrencyRates :copyfrom
INSERT INTO public.currency_rates (id, from_currency, to_currency, rate_date, exchange_rate,
                                   created_at, updated_at, etag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateCurrencyRate :exec
UPDATE public.currency_rates
SET from_currency = $2,
    to_currency   = $3,
    rate_date     = $4,
    exchange_rate = $5,
    updated_at    = $6,
    etag          = $7
WHERE id = $1;

-- name: DeleteCurrencyRate :exec
DELETE
FROM public.currency_rates
WHERE id = $1;

-- name: IsCurrencyRateExists :one
SELECT COUNT(*)
FROM public.currency_rates
WHERE id = ANY ($1::uuid[]);