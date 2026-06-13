DO $$
DECLARE
    added_column TEXT;
BEGIN
    IF to_regclass('public.migration_0006_added_columns') IS NULL THEN
        RETURN;
    END IF;

    FOR added_column IN
        SELECT column_name FROM migration_0006_added_columns
    LOOP
        EXECUTE format('ALTER TABLE users DROP COLUMN IF EXISTS %I', added_column);
    END LOOP;
END $$;

DROP TABLE IF EXISTS migration_0006_added_columns;
