CREATE TABLE IF NOT EXISTS dashboard (
	id VARCHAR(40) PRIMARY KEY,
	year BIGINT NOT NULL,
	month BIGINT NOT NULL,
	request BIGINT NOT NULL,
	staging BIGINT NOT NULL,
	queue BIGINT NOT NULL,
	complete BIGINT NOT NULL,
	failed BIGINT NOT NULL,
	UNIQUE(id),
	UNIQUE(year, month)
);