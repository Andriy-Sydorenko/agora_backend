-- +goose Up
-- Add Google OAuth fields to users table

-- Add google_id (nullable, unique)
ALTER TABLE users
    ADD COLUMN google_id VARCHAR(255);

CREATE UNIQUE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- Add avatar_url (nullable)
ALTER TABLE users
    ADD COLUMN avatar_url VARCHAR(500);

-- Add auth_provider (default 'email')
ALTER TABLE users
    ADD COLUMN auth_provider VARCHAR(20) NOT NULL DEFAULT 'email';

-- Make password nullable (Google OAuth users don't have passwords)
ALTER TABLE users
    ALTER COLUMN password DROP NOT NULL;

-- +goose Down
-- Remove Google OAuth fields

ALTER TABLE users DROP COLUMN IF EXISTS auth_provider;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
DROP INDEX IF EXISTS idx_users_google_id;
ALTER TABLE users DROP COLUMN IF EXISTS google_id;
ALTER TABLE users ALTER COLUMN password SET NOT NULL;
