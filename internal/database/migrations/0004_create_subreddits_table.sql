-- +goose Up
-- Create subreddits table

CREATE TABLE subreddits (
                            id UUID PRIMARY KEY,
                            name VARCHAR(21) NOT NULL UNIQUE,
                            display_name VARCHAR(255) NOT NULL,
                            description VARCHAR(500),
                            icon_url VARCHAR(500),

                            creator_id UUID NOT NULL,

                            member_count INTEGER NOT NULL DEFAULT 0,
                            post_count INTEGER NOT NULL DEFAULT 0,

                            is_public BOOLEAN NOT NULL DEFAULT true,
                            is_nsfw BOOLEAN NOT NULL DEFAULT false,

                            created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
                            updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
                            deleted_at TIMESTAMP WITH TIME ZONE,

                            CONSTRAINT fk_subreddits_creator
                                FOREIGN KEY (creator_id)
                                    REFERENCES users(id)
                                    ON DELETE RESTRICT
);

-- Indexes
CREATE INDEX idx_subreddits_deleted_at ON subreddits(deleted_at);
CREATE INDEX idx_subreddits_creator_id ON subreddits(creator_id);
CREATE INDEX idx_subreddits_name ON subreddits(name);

-- +goose Down
DROP TABLE IF EXISTS subreddits;
