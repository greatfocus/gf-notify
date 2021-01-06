DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM status WHERE name = 'REQUESTED') THEN
        INSERT INTO status (name)
        VALUES
            ('REQUESTED'),
            ('QUEUED'),
            ('PROCESSING'),	
            ('DELIVERED'),
            ('FAILED'),
            ('CANCELLED')
        ON CONFLICT (name) 
        DO NOTHING;
    END IF;
END $$