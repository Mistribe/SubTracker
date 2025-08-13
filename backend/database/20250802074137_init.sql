-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE public.families
(
    id         uuid         NOT NULL PRIMARY KEY,
    name       varchar(100) NOT NULL,
    owner_id   varchar(100) NOT NULL,
    created_at timestamp    NOT NULL,
    updated_at timestamp    NOT NULL,
    etag       varchar(100) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_families_owner
    ON public.families (owner_id);

CREATE TABLE public.family_members
(
    id         uuid         NOT NULL PRIMARY KEY,
    name       varchar(100) NOT NULL,
    family_id  uuid         NOT NULL
        CONSTRAINT fk_families_members
            REFERENCES public.families,
    user_id    varchar(100),
    type       varchar(10)  NOT NULL,
    created_at timestamp    NOT NULL,
    updated_at timestamp    NOT NULL,
    etag       varchar(100) NOT NULL
);

CREATE UNIQUE INDEX family_members_user_id_family_id_uidx
    ON public.family_members (user_id, family_id);

-- Matches the predicate (fm.family_id = p.owner_family_id AND fm.user_id = $1)
CREATE INDEX IF NOT EXISTS idx_family_members_family_user
    ON public.family_members (family_id, user_id);

CREATE TABLE public.labels
(
    id              uuid         NOT NULL PRIMARY KEY,
    owner_type      varchar(10)  NOT NULL,
    owner_family_id uuid
        CONSTRAINT fk_labels_owner_family
            REFERENCES public.families,
    owner_user_id   varchar(50),
    name            varchar(100) NOT NULL,
    key             varchar(100),
    color           varchar(10)  NOT NULL,
    created_at      timestamp    NOT NULL,
    updated_at      timestamp    NOT NULL,
    etag            varchar(100) NOT NULL
);

-- For l.name ILIKE '%term%'
CREATE INDEX IF NOT EXISTS idx_labels_name_trgm
    ON public.labels USING gin (name gin_trgm_ops);

-- FK support for labels.owner_family_id lookups and deletes
CREATE INDEX IF NOT EXISTS idx_labels_owner_family_id
    ON public.labels (owner_family_id);

CREATE TABLE public.providers
(
    id               uuid         NOT NULL PRIMARY KEY,
    owner_type       varchar(20)  NOT NULL,
    owner_family_id  uuid
        CONSTRAINT fk_providers_owner_family
            REFERENCES public.families,
    owner_user_id    varchar(50),
    name             varchar(100) NOT NULL,
    key              varchar(100),
    description      varchar(255),
    icon_url         varchar(255),
    url              varchar(255),
    pricing_page_url varchar(255),
    created_at       timestamp    NOT NULL,
    updated_at       timestamp    NOT NULL,
    etag             varchar(100) NOT NULL
);

-- Personal providers branch
CREATE INDEX IF NOT EXISTS idx_providers_personal_user
    ON public.providers (owner_user_id)
    WHERE owner_type = 'personal';

-- Family providers branch
CREATE INDEX IF NOT EXISTS idx_providers_owner_family_id
    ON public.providers (owner_family_id);

-- For p.name ILIKE '%term%'
CREATE INDEX IF NOT EXISTS idx_providers_name_trgm
    ON public.providers USING gin (name gin_trgm_ops);



CREATE TABLE public.provider_labels
(
    label_id    uuid NOT NULL,
    provider_id uuid NOT NULL,
    CONSTRAINT fk_provider_labels_labels
        FOREIGN KEY (label_id)
            REFERENCES public.labels ON DELETE CASCADE,
    CONSTRAINT fk_provider_labels_providers
        FOREIGN KEY (provider_id)
            REFERENCES public.providers ON DELETE CASCADE,
    PRIMARY KEY (label_id, provider_id)
);

-- Final LEFT JOIN pl ON pl.provider_id = pr.id
CREATE INDEX IF NOT EXISTS idx_provider_labels_provider
    ON public.provider_labels (provider_id);

-- Note: PRIMARY KEY (label_id, provider_id) already covers lookups by label_id.

CREATE TABLE public.provider_plans
(
    id          uuid         NOT NULL PRIMARY KEY,
    name        varchar(100) NOT NULL,
    description varchar(255),
    provider_id uuid         NOT NULL
        CONSTRAINT fk_providers_plans
            REFERENCES public.providers
            ON DELETE CASCADE,
    created_at  timestamp    NOT NULL,
    updated_at  timestamp    NOT NULL,
    etag        varchar(100) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_provider_plans_provider
    ON public.provider_plans (provider_id);


CREATE TABLE public.provider_prices
(
    id         uuid         NOT NULL
        PRIMARY KEY,
    start_date timestamp    NOT NULL,
    end_date   timestamp,
    currency   varchar(3)   NOT NULL,
    amount     numeric      NOT NULL,
    plan_id    uuid         NOT NULL
        CONSTRAINT fk_provider_plans_prices
            REFERENCES public.provider_plans
            ON DELETE CASCADE,
    created_at timestamp    NOT NULL,
    updated_at timestamp    NOT NULL,
    etag       varchar(100) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_provider_prices_plan
    ON public.provider_prices (plan_id);


CREATE TABLE public.subscriptions
(
    id                    uuid         NOT NULL PRIMARY KEY,
    owner_type            varchar(20)  NOT NULL,
    owner_family_id       uuid
        CONSTRAINT fk_subscriptions_owner_family
            REFERENCES public.families,
    owner_user_id         varchar(50),
    friendly_name         varchar(100),
    free_trial_start_date timestamp,
    free_trial_end_date   timestamp,
    provider_id           uuid         NOT NULL
        CONSTRAINT fk_subscriptions_provider
            REFERENCES public.providers,
    plan_id               uuid
        CONSTRAINT fk_subscriptions_plan
            REFERENCES public.provider_plans,
    price_id              uuid
        CONSTRAINT fk_subscriptions_price
            REFERENCES public.provider_prices,
    family_id             uuid
        CONSTRAINT fk_subscriptions_family
            REFERENCES public.families,
    payer_type            varchar(10),
    payer_member_id       uuid
        CONSTRAINT fk_subscriptions_payer_member
            REFERENCES public.family_members,
    start_date            timestamp    NOT NULL,
    end_date              timestamp,
    recurrency            varchar(10)  NOT NULL,
    custom_recurrency     integer,
    custom_price_currency varchar(3),
    custom_price_amount   numeric,
    created_at            timestamp    NOT NULL,
    updated_at            timestamp    NOT NULL,
    etag                  varchar(100) NOT NULL
);

-- FK support and common lookups on subscriptions
CREATE INDEX IF NOT EXISTS idx_subscriptions_provider_id
    ON public.subscriptions (provider_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id
    ON public.subscriptions (plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_price_id
    ON public.subscriptions (price_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_family_id
    ON public.subscriptions (family_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_owner_family_id
    ON public.subscriptions (owner_family_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_payer_member_id
    ON public.subscriptions (payer_member_id);

CREATE TABLE public.subscription_service_users
(
    family_member_id uuid NOT NULL
        CONSTRAINT fk_subscription_service_users_family_member
            REFERENCES public.family_members,
    subscription_id  uuid NOT NULL
        CONSTRAINT fk_subscriptions_service_users
            REFERENCES public.subscriptions ON DELETE CASCADE,
    PRIMARY KEY (family_member_id, subscription_id)
);

CREATE INDEX IF NOT EXISTS idx_subscription_service_users_subscription
    ON public.subscription_service_users (subscription_id);

CREATE TABLE users
(
    id       varchar(100) NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    currency varchar(3)   NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.subscription_service_users;
DROP TABLE public.subscriptions;
DROP TABLE public.provider_prices;
DROP TABLE public.provider_plans;
DROP TABLE public.provider_labels;
DROP TABLE public.providers;
DROP TABLE public.labels;
DROP TABLE public.family_members;
DROP TABLE public.families;
DROP TABLE public.users;
-- +goose StatementEnd
