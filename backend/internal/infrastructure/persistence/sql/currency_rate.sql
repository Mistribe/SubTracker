-- name: GetCurrencyRateById :one
SELECT * FROM public.currency_rates
WHERE id = $1;

-- name: GetCurrencyRateByDate :one
SELECT * FROM public.currency_rates
WHERE from_currency = $1 AND to_currency = $2
  AND rate_date <= $3
ORDER BY rate_date DESC
LIMIT 1;

-- name: GetLatestCurrencyRate :one
SELECT * FROM public.currency_rates
WHERE from_currency = $1 AND to_currency = $2
ORDER BY rate_date DESC
LIMIT 1;

-- name: CreateCurrencyRate :exec
INSERT INTO public.currency_rates (
    id, from_currency, to_currency, rate_date, exchange_rate, 
    created_at, updated_at, etag
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: CreateCurrencyRates :copyfrom
INSERT INTO public.currency_rates (
    id, from_currency, to_currency, rate_date, exchange_rate,
    created_at, updated_at, etag
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8
         );

-- name: UpdateCurrencyRate :exec
UPDATE public.currency_rates
SET from_currency = $2,
    to_currency = $3,
    rate_date = $4,
    exchange_rate = $5,
    updated_at = $6,
    etag = $7
WHERE id = $1;

-- name: DeleteCurrencyRate :exec
DELETE FROM public.currency_rates
WHERE id = $1;

-- name: IsCurrencyRateExists :one
SELECT COUNT(*) FROM public.currency_rates
WHERE id = ANY($1::uuid[]);

-- name: GetCurrencyRates :many
SELECT * FROM public.currency_rates
ORDER BY from_currency, to_currency, rate_date DESC
LIMIT $1 OFFSET $2;

-- name: GetCurrencyRatesCount :one
SELECT COUNT(*) FROM public.currency_rates;