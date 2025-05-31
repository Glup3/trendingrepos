-- +goose Up
ALTER TABLE repositories ADD COLUMN is_archived BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE repositories DROP COLUMN is_archived;
