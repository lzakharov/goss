-- +migrate Up

CREATE TABLE IF NOT EXISTS "user"
(
    id       bigserial not null,
    username text      not null unique,
    password text      not null,
    role     text      not null,

    CONSTRAINT user_pk PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE IF EXISTS "user";
