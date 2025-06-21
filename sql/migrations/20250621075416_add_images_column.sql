-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE images(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id  uuid NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    location TEXT ,
    uploaded_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE images;