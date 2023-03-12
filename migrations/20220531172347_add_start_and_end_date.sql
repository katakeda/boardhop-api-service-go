-- +goose Up
-- +goose StatementBegin
ALTER TABLE "order" DROP COLUMN "created_at";
ALTER TABLE "order" DROP COLUMN "deleted_at";
ALTER TABLE "order" ADD COLUMN "start_date" timestamp;
ALTER TABLE "order" ADD COLUMN "end_date" timestamp;
ALTER TABLE "order" ADD COLUMN "created_at" timestamp NOT NULL DEFAULT NOW();
ALTER TABLE "order" ADD COLUMN "deleted_at" timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "order" DROP COLUMN "start_date";
ALTER TABLE "order" DROP COLUMN "end_date";
-- +goose StatementEnd
