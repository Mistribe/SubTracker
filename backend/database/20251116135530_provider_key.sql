-- +goose Up
-- +goose StatementBegin
ALTER TABLE providers
    ALTER COLUMN key SET NOT NULL;
CREATE UNIQUE INDEX providers_key_uindex
    ON providers (key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX providers_key_uindex;
ALTER TABLE providers
    ALTER COLUMN key DROP NOT NULL;
-- +goose StatementEnd
