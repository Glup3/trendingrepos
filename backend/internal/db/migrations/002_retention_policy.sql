-- +goose Up
SELECT add_retention_policy('stars', INTERVAL '30 days');


-- +goose Down
SELECT remove_retention_policy('stars');
