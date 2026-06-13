DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'users'
          AND column_name = 'password'
    ) THEN
        UPDATE users
        SET password = password_hash
        WHERE password IS NULL;

        ALTER TABLE users ALTER COLUMN password SET NOT NULL;
    END IF;
END $$;
