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
alter table accounts
    alter column currency drop not null;
alter table accounts
    add created_at timestamp default now() not null;
alter table accounts
    add updated_at timestamp default now() not null;
alter table accounts
    add etag varchar(100) default md5(random()::text) not null;

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
        INSERT INTO accounts (id, family_id, role, plan)
        VALUES (NEW.user_id, NEW.family_id, 'user', 'free');
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
-- +goose StatementEnd
