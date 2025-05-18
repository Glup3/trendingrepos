-- +goose Up
CREATE TABLE repositories (
    github_id TEXT NOT NULL PRIMARY KEY,
    name_with_owner TEXT NOT NULL UNIQUE,
    description TEXT,
    stars INTEGER NOT NULL,
    primary_language TEXT
);

CREATE TABLE stars (
    github_id TEXT NOT NULL,
    stars INT NOT NULL,
    time TIMESTAMPTZ NOT NULL
);

SELECT create_hypertable('stars', by_range('time'));



-- +goose Down
DROP TABLE repositories;
DROP TABLE stars;
