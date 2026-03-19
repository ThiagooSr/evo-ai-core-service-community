CREATE TABLE IF NOT EXISTS evo_core_folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL, -- Removed foreign key reference to external accounts table
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_name = 'evo_core_folders'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folders'
        AND column_name = 'name'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'evo_core_folders'
            AND indexname = 'idx_evo_core_folders_name'
        ) THEN
            CREATE INDEX idx_evo_core_folders_name ON evo_core_folders (name);
        END IF;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_name = 'evo_core_folders'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folders'
        AND column_name = 'account_id'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'evo_core_folders'
            AND indexname = 'idx_evo_core_folders_account_id'
        ) THEN
            CREATE INDEX idx_evo_core_folders_account_id ON evo_core_folders (account_id);
        END IF;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_name = 'evo_core_folders'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folders'
        AND column_name = 'account_id'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folders'
        AND column_name = 'name'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'evo_core_folders'
            AND indexname = 'idx_evo_core_folders_account_id_name'
        ) THEN
            CREATE UNIQUE INDEX idx_evo_core_folders_account_id_name ON evo_core_folders (account_id, name);
        END IF;
    END IF;
END
$$;
