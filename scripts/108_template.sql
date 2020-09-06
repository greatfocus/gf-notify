CREATE TABLE IF NOT EXISTS template (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL,
	staticName VARCHAR(20) NOT NULL,
	subject VARCHAR(200) NOT NULL,
	body TEXT NOT NULL,
	paramsCount SMALLINT NOT NULL,
	createdBy BIGINT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedBy BIGINT NULL,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	enabled BOOLEAN NOT NULL default(true),
	deleted BOOLEAN NOT NULL default(false),
	UNIQUE(id),
	UNIQUE(staticName)
);

DO $$ 
BEGIN
	INSERT INTO template (name, staticName, subject, body, paramsCount)
	VALUES
		('OTP', 'otp', 'Please verify your registration', 'Hello, \n\n Thank you for registering and partnering with us. To complete the registration, enter the verification code.\n\n <b>Verification code: $1</b> \n\n Thanks,\n The Respect Team', 1) 
	ON CONFLICT
	DO NOTHING;
END $$;