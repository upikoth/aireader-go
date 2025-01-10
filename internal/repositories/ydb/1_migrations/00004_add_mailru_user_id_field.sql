-- +goose Up
alter table users add column mailru_id text;
-- +goose Down
alter table users drop column mailru_id text;
