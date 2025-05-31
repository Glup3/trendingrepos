-- +goose Up
-- it is somehow possible to have the same repo name for a different github_id
ALTER TABLE repositories DROP CONSTRAINT repositories_name_with_owner_key;

-- +goose Down
ALTER TABLE repositories ADD CONSTRAINT repositories_name_with_owner_key UNIQUE (name_with_owner);
