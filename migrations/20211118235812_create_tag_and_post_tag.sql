-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tag_type" (
    "id" bigserial NOT NULL,
    "name" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "tag" (
    "id" bigserial NOT NULL,
    "type_id" int8 NOT NULL,
    "value" varchar(255) NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_type" FOREIGN KEY ("type_id") REFERENCES "tag_type" ("id")
);

CREATE TABLE "post_tag" (
    "post_id" uuid NOT NULL,
    "tag_id" int8 NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("post_id", "tag_id"),
    CONSTRAINT "fk_post" FOREIGN KEY ("post_id") REFERENCES "post" ("id"),
    CONSTRAINT "fk_tag" FOREIGN KEY ("tag_id") REFERENCES "tag" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "post_tag";
DROP TABLE "tag";
DROP TABLE "tag_type";
-- +goose StatementEnd