CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
DECLARE
	mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
	yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));
BEGIN
	IF NOT EXISTS (SELECT 1 FROM dashboard WHERE year=yr AND month=mnth) THEN
		INSERT INTO dashboard (id, year, month, request, staging, queue, complete, failed)
		VALUES (uuid_generate_v4(), yr, mnth, 0, 0, 0, 0, 0);
	END IF;
END $$;
