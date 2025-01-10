-- +goose Up
create table if not exists users (
     id text,
     email text,
     password_hash text,
     role text,
     primary key (id)
);

create table if not exists sessions (
    id text,
    token text,
    user_id text,
    primary key (id)
);

create table if not exists registrations (
     id text,
     email text,
     confirmation_token text,
     primary key (id)
);

create table if not exists password_recovery_requests (
     id text,
     email text,
     confirmation_token text,
     primary key (id)
);

-- +goose Down
drop table if exists users;
drop table if exists sessions;
drop table if exists registrations;
drop table if exists password_recovery_requests;
