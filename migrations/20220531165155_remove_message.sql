-- +goose Up
-- +goose StatementBegin
ALTER TABLE "order" DROP COLUMN "message";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "order" ADD COLUMN "message" varchar(255);
-- +goose StatementEnd
