-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.currency_rates
(
    id            uuid         NOT NULL PRIMARY KEY,
    from_currency varchar(3)   NOT NULL,
    to_currency   varchar(3)   NOT NULL,
    rate_date     date         NOT NULL,
    exchange_rate numeric      NOT NULL,
    created_at    timestamp    NOT NULL,
    updated_at    timestamp    NOT NULL,
    etag          varchar(100) NOT NULL
);

-- Create an index on the currency pair and date for efficient lookups
CREATE INDEX currency_rates_from_to_date_idx
    ON public.currency_rates (from_currency, to_currency, rate_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.currency_rates;
-- +goose StatementEnd