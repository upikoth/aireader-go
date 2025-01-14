-- +goose Up
create table if not exists voices (
    name text,
    source text,
    primary key (name)
);

-- +goose Down
drop table if exists voices;
