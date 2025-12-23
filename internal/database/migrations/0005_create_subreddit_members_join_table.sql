-- +goose Up
-- Create subreddit_members join table for many-to-many relationship

CREATE TABLE subreddit_members (
                                   subreddit_id UUID NOT NULL,
                                   user_id UUID NOT NULL,

                                   PRIMARY KEY (subreddit_id, user_id),

                                   CONSTRAINT fk_subreddit_members_subreddit
                                       FOREIGN KEY (subreddit_id)
                                           REFERENCES subreddits(id)
                                           ON DELETE CASCADE,

                                   CONSTRAINT fk_subreddit_members_user
                                       FOREIGN KEY (user_id)
                                           REFERENCES users(id)
                                           ON DELETE CASCADE
);

-- Index for querying user's subreddits
CREATE INDEX idx_subreddit_members_user_id ON subreddit_members(user_id);

-- +goose Down
DROP TABLE IF EXISTS subreddit_members;
