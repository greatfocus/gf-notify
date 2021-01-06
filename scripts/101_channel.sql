CREATE TABLE IF NOT EXISTS channel (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL,
	staticName VARCHAR(20) NOT NULL,
	priority SMALLINT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	enabled BOOLEAN NOT NULL default(true),
	UNIQUE(id),
	UNIQUE(staticName)
);