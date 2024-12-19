CREATE TABLE artists
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE release_artists
(
    id         SERIAL PRIMARY KEY,
    release_id INTEGER REFERENCES releases (id),
    artist_id  INTEGER REFERENCES artists (id)
);
