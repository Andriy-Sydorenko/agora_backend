-- +goose Up
-- Ensure username and email have proper size constraints

-- Update username to VARCHAR(255)
ALTER TABLE users
ALTER COLUMN username TYPE VARCHAR(255);

-- Update email to VARCHAR(100) if it's different
ALTER TABLE users
ALTER COLUMN email TYPE VARCHAR(255);

-- Ensure auth_provider is VARCHAR(20) (should already be set from 0002, but ensuring consistency)
ALTER TABLE users
ALTER COLUMN auth_provider TYPE VARCHAR(20);

-- +goose Down
-- No rollback needed - these are just type constraints
-- Original migrations already have correct types
