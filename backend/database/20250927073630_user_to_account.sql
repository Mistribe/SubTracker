-- +goose Up
-- +goose StatementBegin
ALTER INDEX users_pk RENAME TO accounts_pk;

ALTER TABLE users
    RENAME TO accounts;

ALTER TABLE accounts
    ADD family_id uuid
        CONSTRAINT accounts_families_id_fk
            REFERENCES families;
ALTER TABLE accounts
    ADD role varchar(10) NOT NULL DEFAULT 'user';
ALTER TABLE subscriptions
    DROP COLUMN plan_id;
ALTER TABLE subscriptions
    DROP COLUMN price_id;
DROP TABLE provider_prices;
DROP TABLE provider_plans;
ALTER TABLE subscription_service_users
    RENAME TO subscription_family_users;

CREATE OR REPLACE FUNCTION update_account_family_id()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.user_id IS NULL THEN
        UPDATE accounts
        SET family_id = NULL
        WHERE id = OLD.user_id;
        RETURN NEW;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM accounts a WHERE a.id = NEW.user_id) THEN
        INSERT INTO accounts (id, family_id, role, plan, currency)
        VALUES (NEW.user_id, NEW.family_id, 'user', 'free', '');
    ELSE
        UPDATE accounts
        SET family_id = NEW.family_id
        WHERE id = NEW.user_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_account_family_id_on_member_add
    AFTER INSERT
    ON family_members
    FOR EACH ROW
EXECUTE FUNCTION update_account_family_id();

CREATE TRIGGER update_account_family_id_on_member_update
    AFTER UPDATE OF family_id
    ON family_members
    FOR EACH ROW
EXECUTE FUNCTION update_account_family_id();


CREATE OR REPLACE FUNCTION remove_account_family_id()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE accounts
    SET family_id = NULL
    WHERE id = OLD.user_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER remove_account_family_id_on_member_delete
    AFTER DELETE
    ON family_members
    FOR EACH ROW
EXECUTE FUNCTION remove_account_family_id();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accounts
    DROP CONSTRAINT accounts_families_id_fk;

ALTER TABLE accounts
    DROP COLUMN family_id;

ALTER TABLE accounts
    RENAME TO users;

ALTER INDEX accounts_pk RENAME TO users_pk;

ALTER TABLE accounts
    DROP COLUMN role;

CREATE TABLE provider_plans
(
    id          uuid PRIMARY KEY,
    provider_id varchar(255) NOT NULL,
    name        varchar(255) NOT NULL,
    description text,
    created_at  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE provider_prices
(
    id         uuid PRIMARY KEY,
    plan_id    uuid         NOT NULL REFERENCES provider_plans (id),
    price_id   varchar(255) NOT NULL,
    amount     decimal      NOT NULL,
    currency   varchar(3)   NOT NULL,
    interval   varchar(10)  NOT NULL,
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE subscriptions
    ADD COLUMN plan_id  uuid REFERENCES provider_plans (id),
    ADD COLUMN price_id uuid REFERENCES provider_prices (id);

ALTER TABLE subscription_family_users
    RENAME TO subscription_service_users;

DROP TRIGGER IF EXISTS set_account_family_id_on_member_add ON family_members;
DROP FUNCTION IF EXISTS update_account_family_id();
-- +goose StatementEnd
