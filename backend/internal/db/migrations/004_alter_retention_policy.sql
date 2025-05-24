-- +goose Up
SELECT remove_retention_policy('stars');
SELECT add_retention_policy('stars', INTERVAL '35 days');

-- +goose Down
SELECT remove_retention_policy('stars');
SELECT add_retention_policy('stars', INTERVAL '30 days');
