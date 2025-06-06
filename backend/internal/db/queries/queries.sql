-- name: CreateTempRepositories :exec
CREATE TEMPORARY TABLE temp_repositories (LIKE repositories INCLUDING ALL) ON COMMIT DROP;

-- name: InsertTempRepositories :copyfrom
INSERT INTO temp_repositories (github_id, name_with_owner, description, stars, primary_language, is_archived)
  VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertRepositories :exec
INSERT INTO repositories
SELECT * FROM temp_repositories
ON CONFLICT (github_id) DO UPDATE
SET
  stars = EXCLUDED.stars,
  description = EXCLUDED.description,
  primary_language = EXCLUDED.primary_language,
  is_archived = EXCLUDED.is_archived;

-- name: InsertStars :exec
INSERT INTO stars (github_id, stars, time)
SELECT github_id, stars, NOW() FROM temp_repositories;

