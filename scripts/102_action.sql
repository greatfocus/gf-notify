CREATE TABLE IF NOT EXISTS actions (
	id SERIAL PRIMARY KEY,
	name VARCHAR(30) NOT NULL,
	description VARCHAR(100) NOT NULL,
	createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(name)
);

INSERT INTO actions (name, description, deleted, enabled)
VALUES
	('user_edit', 'Action to allow user edit', false, true),
	('user_delete', 'Action to allow user delete', false, true),
	('user_approve', 'Action to allow user approve', false, true),
	('user_reject', 'Action to allow user reject', false, true),
	('user_activate', 'Action to allow user activate', false, true),
	('user_deactivate', 'Action to allow user deactivate', false, true)
ON CONFLICT (name) 
DO NOTHING;
