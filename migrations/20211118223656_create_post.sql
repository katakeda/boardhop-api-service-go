-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS rate;

CREATE TYPE rate AS ENUM ('hour', 'day', 'week', 'month');

CREATE TABLE "post" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL,
    "title" varchar(255) NOT NULL,
    "description" text,
    "price" float4 NOT NULL,
    "rate" rate NOT NULL,
    "pickup_latitude" float8,
    "pickup_longitude" float8,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "deleted_at" timestamp,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "post";
-- +goose StatementEnd
