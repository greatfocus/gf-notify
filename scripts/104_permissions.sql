CREATE TABLE IF NOT EXISTS permissions (
	id SERIAL PRIMARY KEY,
	roleid INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	actionid INTEGER NOT NULL REFERENCES actions(id) ON DELETE CASCADE,
	createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false)
);

DO $$ 
DECLARE
	adminroleid INTEGER := (select id from roles where name='admin');
	staffroleid INTEGER := (select id from roles where name='staff');
	agentroleid INTEGER := (select id from roles where name='agent');
	customerroleid INTEGER := (select id from roles where name='customer');


	user_edit INTEGER := (select id from actions where name='user_edit');
	user_delete INTEGER := (select id from actions where name='user_delete');
	user_approve INTEGER := (select id from actions where name='user_approve');
	user_reject INTEGER := (select id from actions where name='user_reject');
	user_activate INTEGER := (select id from actions where name='user_activate');
	user_deactivate INTEGER := (select id from actions where name='user_deactivate');
	
BEGIN 
	INSERT INTO permissions (roleid, actionid, deleted, enabled)
	VALUES
		-- admin
		(adminroleid, user_edit, false, true),
		(adminroleid, user_delete, false, true),
		(adminroleid, user_approve, false, true),
		(adminroleid, user_reject, false, true),
		(adminroleid, user_activate, false, true),
		(adminroleid, user_deactivate, false, true),

		-- staff
		(staffroleid, user_edit, false, true),
		(staffroleid, user_delete, false, true),
		(staffroleid, user_activate, false, true),
		(staffroleid, user_deactivate, false, true)

		-- agent
		-- customer
	ON CONFLICT
	DO NOTHING;
END $$;
