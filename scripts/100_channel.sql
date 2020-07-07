CREATE TABLE IF NOT EXISTS channel (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL,
	staticName VARCHAR(20) NOT NULL,
	priority SMALLINT NULL,
	createdBy BIGINT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedBy BIGINT NULL,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	enabled BOOLEAN NOT NULL default(true),
	UNIQUE(id),
	UNIQUE(staticName)
);


INSERT INTO channel (name, staticName, priority)
VALUES
	('sms', 'sms', 1),
	('email', 'email', 2)
ON CONFLICT (staticName) 
DO NOTHING;