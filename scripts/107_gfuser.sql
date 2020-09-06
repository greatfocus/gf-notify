CREATE TABLE IF NOT EXISTS gfuser (
	id BIGSERIAL PRIMARY KEY,
	relatedId INTEGER NOT NULL,
	email VARCHAR(50) NOT NULL,
	key VARCHAR(100) NOT NULL,
	createdBy BIGINT NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedBy BIGINT NOT NULL,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	enabled BOOLEAN NOT NULL default(true),
	deleted BOOLEAN NOT NULL default(false),
	UNIQUE(email),
	UNIQUE(email, key)
);


DO $$ 
BEGIN
	INSERT INTO gfuser (relatedId, email, key, createdBy, updatedBy)
	VALUES
		(1, 'mucunga90@gmail.com', '$2a$04$z7wANFLrGo7c/vr0gG22FOF7adJDDI/Xto7/Mjotirs3SR62Xjh7e', 1, 1) 
	ON CONFLICT
	DO NOTHING;
END $$;
