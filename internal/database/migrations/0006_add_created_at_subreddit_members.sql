-- +goose Up
ALTER TABLE subreddit_members
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE subreddit_members
    DROP COLUMN created_at;
