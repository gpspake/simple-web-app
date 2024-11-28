CREATE TABLE artists
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT    NOT NULL
);

CREATE TABLE release_artists
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    release_id INTEGER REFERENCES artists(id),
    artist_id  INTEGER REFERENCES releases(id)
);
