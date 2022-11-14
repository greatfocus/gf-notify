CREATE TABLE IF NOT EXISTS channel (
	id VARCHAR(40) PRIMARY KEY,
	name TEXT NOT NULL,
	key TEXT NOT NULL,
	url TEXT NOT NULL,
	username TEXT NOT NULL,
	pass TEXT NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	enabled BOOLEAN NOT NULL default(true),
	UNIQUE(id),
	UNIQUE(key)
);