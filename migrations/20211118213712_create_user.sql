-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "email" text NOT NULL UNIQUE,
    "first_name" varchar(50) NOT NULL,
    "last_name" varchar(50) NOT NULL,
    "phone" varchar(50),
    "avatar_url" text,
    "google_auth_id" text,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd
