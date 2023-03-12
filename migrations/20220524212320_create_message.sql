-- +goose Up
-- +goose StatementBegin
CREATE TABLE "message" (
    "id" bigserial NOT NULL,
    "user_id" uuid NOT NULL,
    "post_id" uuid,
    "order_id" uuid,
    "message" varchar(255),
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "message";
-- +goose StatementEnd
