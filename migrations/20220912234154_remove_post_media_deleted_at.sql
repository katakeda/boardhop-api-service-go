-- +goose Up
-- +goose StatementBegin
ALTER TABLE "post_media" DROP COLUMN "deleted_at";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "post_media" ADD COLUMN "deleted_at" timestamp;
-- +goose StatementEnd
