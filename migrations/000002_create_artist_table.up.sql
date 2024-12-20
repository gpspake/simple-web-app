CREATE TABLE artist
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE release_artist
(
    id         SERIAL PRIMARY KEY,
    release_id INTEGER REFERENCES release (id),
    artist_id  INTEGER REFERENCES artist (id)
);
