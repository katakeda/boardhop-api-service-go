-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS media_type;

CREATE TYPE media_type AS ENUM ('image', 'video');

CREATE TABLE "post_media" (
    "id" bigserial NOT NULL,
    "post_id" uuid NOT NULL,
    "media_url" text NOT NULL,
    "type" media_type NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "deleted_at" timestamp,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_post" FOREIGN KEY ("post_id") REFERENCES "post" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "post_media";
-- +goose StatementEnd
