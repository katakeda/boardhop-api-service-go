-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE "category" (
    "id" bigserial NOT NULL,
    "parent_id" int8,
    "path" LTREE,
    "value" varchar(255) NOT NULL,
    "label" varchar(255) NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_parent" FOREIGN KEY ("parent_id") REFERENCES "category" ("id")
);

CREATE INDEX "category_path_idx" ON "category" USING GIST ("path");

CREATE OR REPLACE FUNCTION update_category_path() RETURNS TRIGGER AS $$
    DECLARE
        NEW_PATH ltree;
    BEGIN
        IF NEW.parent_id IS NULL THEN
            NEW.path = 'root'::ltree;
        ELSEIF TG_OP = 'INSERT' OR OLD.parent_id IS NULL OR OLD.parent_id != NEW.parent_id THEN
            SELECT path || id::text FROM category WHERE id = NEW.parent_id INTO NEW_PATH;
            IF NEW_PATH IS NULL THEN
                RAISE EXCEPTION 'Invalid parent_id %', NEW.parent_id;
            END IF;
            NEW.path = NEW_PATH;
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER category_path_trigger
    BEFORE INSERT OR UPDATE ON category
    FOR EACH ROW EXECUTE PROCEDURE update_category_path();

CREATE TABLE "post_category" (
    "post_id" uuid NOT NULL,
    "category_id" int8 NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("post_id", "category_id"),
    CONSTRAINT "fk_post" FOREIGN KEY ("post_id") REFERENCES "post" ("id"),
    CONSTRAINT "fk_category" FOREIGN KEY ("category_id") REFERENCES "category" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS "category_path_trigger" ON "boardhop"."category";
DROP TABLE "post_category";
DROP TABLE "category";
-- +goose StatementEnd
