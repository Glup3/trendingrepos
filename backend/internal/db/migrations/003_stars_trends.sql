-- +goose NO TRANSACTION
-- +goose Up

CREATE MATERIALIZED VIEW stars_daily
WITH
	(timescaledb.continuous) AS
SELECT
	time_bucket ('1 day', time) AS bucket,
	github_id,
	FIRST (stars, time) AS stars_earliest,
	LAST (stars, time) AS stars_latest
FROM
	stars
GROUP BY
	time_bucket ('1 day', time),
	github_id
WITH
	DATA;


CREATE MATERIALIZED VIEW stars_trend_monthly AS
SELECT
	r.github_id,
	r.name_with_owner,
	r.primary_language,
	r.description,
	COALESCE(today.stars_latest, 0) AS stars_today,
	COALESCE(past.stars_earliest, 0) AS stars_prev,
	COALESCE(today.stars_latest, 0) - COALESCE(past.stars_earliest, 0) AS stars_diff
FROM
	stars_daily AS today
	LEFT JOIN stars_daily AS past ON today.github_id = past.github_id
	AND past.bucket = time_bucket ('1 day', NOW() - INTERVAL '30 days')
	LEFT JOIN repositories AS r ON today.github_id = r.github_id
WHERE
	today.bucket = time_bucket ('1 day', NOW())
GROUP BY
	r.github_id,
	today.stars_latest,
	past.stars_earliest;


CREATE MATERIALIZED VIEW stars_trend_daily AS
SELECT
	r.github_id,
	r.name_with_owner,
	r.primary_language,
	r.description,
	COALESCE(today.stars_latest, 0) AS stars_today,
	COALESCE(past.stars_earliest, 0) AS stars_prev,
	COALESCE(today.stars_latest, 0) - COALESCE(past.stars_earliest, 0) AS stars_diff
FROM
	stars_daily AS today
	LEFT JOIN stars_daily AS past ON today.github_id = past.github_id
	AND past.bucket = time_bucket ('1 day', NOW() - INTERVAL '1 days')
	LEFT JOIN repositories AS r ON today.github_id = r.github_id
WHERE
	today.bucket = time_bucket ('1 day', NOW())
GROUP BY
	r.github_id,
	today.stars_latest,
	past.stars_earliest;


CREATE MATERIALIZED VIEW stars_trend_weekly AS
SELECT
	r.github_id,
	r.name_with_owner,
	r.primary_language,
	r.description,
	COALESCE(today.stars_latest, 0) AS stars_today,
	COALESCE(past.stars_earliest, 0) AS stars_prev,
	COALESCE(today.stars_latest, 0) - COALESCE(past.stars_earliest, 0) AS stars_diff
FROM
	stars_daily AS today
	LEFT JOIN stars_daily AS past ON today.github_id = past.github_id
	AND past.bucket = time_bucket ('1 day', NOW() - INTERVAL '7 days')
	LEFT JOIN repositories AS r ON today.github_id = r.github_id
WHERE
	today.bucket = time_bucket ('1 day', NOW())
GROUP BY
	r.github_id,
	today.stars_latest,
	past.stars_earliest;


-- +goose Down
DROP MATERIALIZED VIEW stars_daily CASCADE;
