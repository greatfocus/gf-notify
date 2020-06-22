DO $$ 
DECLARE
	mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
	yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));


BEGIN
	EXECUTE format('
	CREATE TABLE IF NOT EXISTS messageOut%s%s (
		id BIGSERIAL PRIMARY KEY,
		channel VARCHAR(10) NOT NULL,
		recipient VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		createdBy BIGINT NOT NULL,
		createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		statusId INTEGER NOT NULL,
		attempts SMALLINT NOT NULL,
		priority SMALLINT NOT NULL,
		refId BIGINT NOT NULL,
		UNIQUE(id, createdOn)
	);', yr, mnth);
END $$;
