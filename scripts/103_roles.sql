CREATE TABLE IF NOT EXISTS roles (
	id SERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL,
	description VARCHAR(100) NOT NULL,
	createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(name)
);

INSERT INTO roles (name, description, deleted, enabled)
VALUES
	('admin', 'Role for admin', false, true),
	('staff', 'Role for staff', false, true),
	('agent', 'Role for agent', false, true),	
	('default', 'Role for default customer', false, true)
ON CONFLICT (name) 
DO NOTHING;
