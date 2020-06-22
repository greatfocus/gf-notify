CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	type VARCHAR(10) NOT NULL,
	firstname VARCHAR(20) NOT NULL,
	middlename VARCHAR(20) NOT NULL,
	lastname VARCHAR(20) NOT NULL,
	mobilenumber VARCHAR(14) NOT NULL,
	email VARCHAR(50) NOT NULL,
	password VARCHAR(100) NOT NULL,
	failedattempts FLOAT default (0),
	lastattempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	lastchange TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expireddate TIMESTAMP NOT NULL,
	createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(20) NOT NULL,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(mobilenumber),
	UNIQUE(email),
	UNIQUE(email, mobilenumber)
);

INSERT INTO users (type, firstname, middlename, lastname, mobilenumber, email, password, expireddate, status, enabled)
VALUES
	('password', 'admin', 'admin', 'admin', '0780904371', 'mucunga90@gmail.com', '$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe', '2025-06-01 08:22:17.460493', 'USER.APPROVED', true) 
ON CONFLICT (email) 
DO NOTHING;
