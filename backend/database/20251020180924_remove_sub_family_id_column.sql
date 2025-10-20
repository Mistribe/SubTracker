-- +goose Up
-- +goose StatementBegin
alter table subscriptions
    drop column family_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE subscriptions ADD COLUMN family_id uuid;
ALTER TABLE subscriptions ADD CONSTRAINT subscriptions_family_id_fkey FOREIGN KEY (family_id) REFERENCES families(id);
-- +goose StatementEnd
