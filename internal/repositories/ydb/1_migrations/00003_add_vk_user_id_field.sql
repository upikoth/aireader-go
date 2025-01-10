-- +goose Up
alter table users add column vk_id text;
-- +goose Down
alter table users drop column vk_id text;
