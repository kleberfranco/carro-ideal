DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'users'
          AND column_name = 'password'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
    END IF;
END $$;
