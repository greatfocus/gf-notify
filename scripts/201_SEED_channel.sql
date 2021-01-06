DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM channel WHERE staticName = 'email') THEN
        INSERT INTO channel (name, staticName, priority)
        VALUES
            ('sms', 'sms', 1),
            ('email', 'email', 2)
        ON CONFLICT (staticName) 
        DO NOTHING;
    END IF;
END $$