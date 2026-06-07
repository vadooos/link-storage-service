CREATE TABLE IF NOT EXISTS links
(
    id           BIGSERIAL PRIMARY KEY,
    short_code   TEXT        NOT NULL UNIQUE,
    original_url TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    visits       BIGINT      NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_links_short_code ON links (short_code);