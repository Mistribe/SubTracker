-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD plan varchar(30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN plan;
-- +goose StatementEnd
