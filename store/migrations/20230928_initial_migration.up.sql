-- +goose Up
CREATE TABLE users (
    id uuid NOT NULL,
    firstname text NOT NULL,
    lastname text NOT NULL,
    nickname text NOT NULL,
    password  text NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp default current_timestamp,
    CONSTRAINT "pk_user_id" PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;