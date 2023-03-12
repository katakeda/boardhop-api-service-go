-- +goose Up
-- +goose StatementBegin
INSERT INTO "category" ("value", "label") VALUES ('surfboard', 'サーフボード');
INSERT INTO "category" ("value", "label") VALUES ('snowboard', 'スノーボード');
INSERT INTO "category" ("parent_id", "value", "label") VALUES ((SELECT "id" FROM "category" WHERE "value" = 'surfboard'), 'shortboard', 'ショートボード');
INSERT INTO "category" ("parent_id", "value", "label") VALUES ((SELECT "id" FROM "category" WHERE "value" = 'surfboard'), 'longboard', 'ロングボード');
INSERT INTO "category" ("parent_id", "value", "label") VALUES ((SELECT "id" FROM "category" WHERE "value" = 'shortboard'), 'eps', 'EPS');
INSERT INTO "category" ("parent_id", "value", "label") VALUES ((SELECT "id" FROM "category" WHERE "value" = 'shortboard'), 'pu', 'PU');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Skill Level', 'beginner', '初心者');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Skill Level', 'intermediate', '中級者');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Skill Level', 'advanced', '上級者');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Surfboard Brand', 'ci', 'CI');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Surfboard Brand', 'pyzel', 'Pyzel');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Surfboard Brand', 'js', 'JS');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Surfboard Brand', 'mayhem', 'Mayhem');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Surfboard Brand', 'hs', 'Hayden Shapes');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Surfboard Brand', 'firewire', 'Firewire');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Snowboard Brand', 'burton', 'Burton');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Snowboard Brand', 'arbor', 'Arbor');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Snowboard Brand', 'union', 'Union');
INSERT INTO "tag" ("type", "value", "label") VALUES ('Snowboard Brand', 'k2', 'K2');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "category" WHERE "value" = 'pu';
DELETE FROM "category" WHERE "value" = 'eps';
DELETE FROM "category" WHERE "value" = 'longboard';
DELETE FROM "category" WHERE "value" = 'shortboard';
DELETE FROM "category" WHERE "value" = 'snowboard';
DELETE FROM "category" WHERE "value" = 'surfboard';
DELETE FROM "tag" WHERE "value" = 'beginner';
DELETE FROM "tag" WHERE "value" = 'intermediate';
DELETE FROM "tag" WHERE "value" = 'advanced';
DELETE FROM "tag" WHERE "value" = 'ci';
DELETE FROM "tag" WHERE "value" = 'pyzel';
DELETE FROM "tag" WHERE "value" = 'js';
DELETE FROM "tag" WHERE "value" = 'mayhem';
DELETE FROM "tag" WHERE "value" = 'hs';
DELETE FROM "tag" WHERE "value" = 'firewire';
DELETE FROM "tag" WHERE "value" = 'burton';
DELETE FROM "tag" WHERE "value" = 'arbor';
DELETE FROM "tag" WHERE "value" = 'union';
DELETE FROM "tag" WHERE "value" = 'k2';
-- +goose StatementEnd
