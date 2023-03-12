-- +goose Up
-- +goose StatementBegin
INSERT INTO "user" ("email", "first_name", "last_name", "phone") VALUES ('test@theboardhop.com', 'Test', 'Account', '0000000000');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "user" WHERE "email" = 'test@theboardhop.com';
-- +goose StatementEnd
