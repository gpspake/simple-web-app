CREATE TABLE releases_fts
(
    id              SERIAL PRIMARY KEY,
    release_id      INTEGER NOT NULL,
    release_name    TEXT    NOT NULL,
    release_year    INTEGER NOT NULL,
    artist_name     TEXT    NOT NULL,
    tsvector_column TSVECTOR
);

-- Populate the tsvector column
UPDATE releases_fts
SET tsvector_column = to_tsvector(release_name || ' ' || artist_name || ' ' || release_year::TEXT);

-- Add an index for full-text search
CREATE INDEX idx_releases_fts_tsvector ON releases_fts USING gin (tsvector_column);
