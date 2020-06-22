CREATE TABLE IF NOT EXISTS rights (
	id SERIAL PRIMARY KEY,
	roleid INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	userid INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(20) NOT NULL,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(roleid, userid)
);

DO $$ 
DECLARE
	roleid INTEGER := (select id from roles where name='admin');
	userid INTEGER := (select id from users where email='mucunga90@gmail.com');
BEGIN 
	INSERT INTO rights (roleid, userid, status, deleted, enabled)
	VALUES
		(roleid, userid, 'RIGHT.APPROVED', false, true)
	ON CONFLICT
	DO NOTHING;
END $$;
