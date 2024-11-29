-- Create full text search table for releases
CREATE VIRTUAL TABLE releases_fts USING fts5
(
    release_id UNINDEXED,
    release_name,
    release_year,
    artist_name,
    tokenize="trigram"
);

-- Populate full text search table for releases
-- INSERT INTO releases_fts (release_id, artist_name, release_name, release_year)
-- SELECT
--     releases.id AS release_id,
--     artists.name AS artist_name,
--     releases.name AS release_name,
--     releases.year AS release_year
-- FROM
--     release_artists
--         JOIN
--     artists ON release_artists.artist_id = artists.id
--         JOIN
--     releases ON release_artists.release_id = releases.id;

-- Trigger to update full text search table after release inserts
CREATE TRIGGER releases_ai AFTER INSERT ON releases
BEGIN
    INSERT INTO releases_fts (release_id, release_name, release_year, artist_name)
    SELECT
        releases.id AS release_id,
        releases.name AS release_name,
        releases.year AS release_year,
        artists.name AS artist_name
    FROM
        release_artists
            JOIN
        artists ON release_artists.artist_id = artists.id
            JOIN
        releases ON release_artists.release_id = releases.id;
END;

-- Trigger to update full text search table after release updates
CREATE TRIGGER releases_au AFTER UPDATE ON releases
BEGIN
    UPDATE releases_fts
    SET release_name = NEW.name,
        release_year = NEW.year
    WHERE release_name = OLD.name AND release_year = OLD.year;
END;

-- Trigger to update full text search table after release deletes
CREATE TRIGGER releases_ad AFTER DELETE ON releases
BEGIN
    DELETE FROM releases_fts
    WHERE release_id = OLD.id;
END;

-- Todo: Trigger to update full text search table after release_artist inserts
-- CREATE TRIGGER release_artists_ai AFTER INSERT ON release_artists
-- END;

-- Todo: Trigger to update full text search table after release_artist updates
-- CREATE TRIGGER release_artists_au AFTER UPDATE ON release_artists
-- END;

-- Todo: Trigger to update full text search table after release_artist deletes
-- CREATE TRIGGER release_artists_au AFTER DELETE ON release_artists
-- END;

-- Trigger to update full text search table after artist updates
CREATE TRIGGER artists_au AFTER UPDATE ON artists
BEGIN
    UPDATE releases_fts
    SET artist_name = NEW.name
    WHERE artist_name = OLD.name;
END;
