DO $$ 
DECLARE
	mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
	yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));


BEGIN
	EXECUTE format('
	CREATE TABLE IF NOT EXISTS failed%s%s (
		id BIGSERIAL,
		channelId INTEGER REFERENCES channel(id),
		recipient VARCHAR(100) NOT NULL,
		subject VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		createdBy BIGINT NOT NULL,
		createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expireOn TIMESTAMP NOT NULL,
		statusId INTEGER REFERENCES status(id),
		attempts SMALLINT NOT NULL,
		priority SMALLINT NOT NULL,
		reference TEXT NOT NULL,
		UNIQUE(id),		
		UNIQUE(id, createdOn)
	);', yr, mnth);
END $$;
