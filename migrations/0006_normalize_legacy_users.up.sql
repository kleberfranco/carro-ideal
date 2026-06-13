CREATE TABLE IF NOT EXISTS migration_0006_added_columns (
    column_name TEXT PRIMARY KEY
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'users'
          AND column_name = 'password_hash'
    ) THEN
        ALTER TABLE users ADD COLUMN password_hash TEXT;
        INSERT INTO migration_0006_added_columns (column_name) VALUES ('password_hash');

        IF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'users'
              AND column_name = 'password'
        ) THEN
            UPDATE users SET password_hash = password WHERE password_hash IS NULL;
        END IF;

        ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'users'
          AND column_name = 'role'
    ) THEN
        ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';
        INSERT INTO migration_0006_added_columns (column_name) VALUES ('role');
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'users'
          AND column_name = 'updated_at'
    ) THEN
        ALTER TABLE users ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();
        INSERT INTO migration_0006_added_columns (column_name) VALUES ('updated_at');
    END IF;
END $$;
