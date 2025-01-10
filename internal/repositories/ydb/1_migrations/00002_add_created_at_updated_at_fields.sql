-- +goose Up
alter table sessions add column created_at timestamp;

alter table registrations add column created_at timestamp;

alter table password_recovery_requests add column created_at timestamp;

alter table users add column created_at timestamp;
alter table users add column updated_at timestamp;
-- +goose Down
alter table sessions drop column created_at;

alter table registrations drop column created_at timestamp;

alter table password_recovery_requests drop column created_at timestamp;

alter table users drop column created_at timestamp;
alter table users drop column updated_at timestamp;
