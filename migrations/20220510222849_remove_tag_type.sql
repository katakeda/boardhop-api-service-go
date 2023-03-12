-- +goose Up
-- +goose StatementBegin
ALTER TABLE "tag" DROP CONSTRAINT "fk_type";
ALTER TABLE "tag" ALTER COLUMN "type_id" SET DATA TYPE varchar(255);
ALTER TABLE "tag" RENAME COLUMN "type_id" TO "type";
ALTER TABLE "tag" ADD COLUMN "label" varchar(255) NOT NULL;
DROP TABLE "tag_type";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE "tag_type" (
    "id" bigserial NOT NULL,
    "name" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
);
ALTER TABLE "tag" DROP COLUMN "label";
ALTER TABLE "tag" RENAME COLUMN "type" TO "type_id";
ALTER TABLE "tag" ALTER COLUMN "type_id" SET DATA TYPE int8;
ALTER TABLE "tag" CONSTRAINT "fk_type" FOREIGN KEY ("type_id") REFERENCES "tag_type" ("id");
-- +goose StatementEnd
