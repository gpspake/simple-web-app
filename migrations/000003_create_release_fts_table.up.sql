CREATE TABLE release_fts
(
    id              SERIAL PRIMARY KEY,
    release_id      INTEGER NOT NULL,
    release_title    TEXT    NOT NULL,
    release_year    INTEGER NOT NULL,
    artist_name     TEXT    NOT NULL,
    tsvector_column TSVECTOR
);

-- Populate the tsvector column
UPDATE release_fts
SET tsvector_column = to_tsvector(release_title || ' ' || artist_name || ' ' || release_year::TEXT);

-- Add an index for full-text search
CREATE INDEX idx_release_fts_tsvector ON release_fts USING gin (tsvector_column);
