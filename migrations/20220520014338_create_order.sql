-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS order_status;

CREATE TYPE order_status AS ENUM ('pending', 'complete', 'canceled');

CREATE TABLE "order" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "post_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "payment_id" varchar(255) NOT NULL,
    "status" order_status NOT NULL,
    "message" varchar(255),
    "quantity" int4 NOT NULL,
    "total" float4 NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "deleted_at" timestamp,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_post" FOREIGN KEY ("post_id") REFERENCES "post" ("id"),
    CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "order";
-- +goose StatementEnd
