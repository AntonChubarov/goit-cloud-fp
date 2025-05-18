CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE links (
                       id            UUID PRIMARY KEY        DEFAULT uuid_generate_v4(),
                       short_code    VARCHAR(12) UNIQUE NOT NULL,
                       original_url  TEXT           NOT NULL,
                       created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
                       clicks        BIGINT         NOT NULL DEFAULT 0
);

-- indexes for fast look-up
CREATE UNIQUE INDEX idx_links_short_code ON links(short_code);
