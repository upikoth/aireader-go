-- +goose Up
alter table users add column yandex_id text;
-- +goose Down
alter table users drop column yandex_id text;
