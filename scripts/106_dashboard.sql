DO $$ 
DECLARE
	mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
	yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));

BEGIN
	CREATE TABLE IF NOT EXISTS dashboard (
		year INTEGER NOT NULL,
		month INTEGER NOT NULL,
		staging INTEGER NOT NULL,
		queue INTEGER NOT NULL,
		complete INTEGER NOT NULL,
		failed INTEGER NOT NULL
	);

	IF (SELECT count(staging) FROM dashboard WHERE year=yr AND month=mnth) < 1 THEN
		INSERT INTO dashboard (year, month, staging, queue, complete, failed)
		VALUES (yr, mnth, 0, 0, 0, 0);
	END IF;
END $$;
