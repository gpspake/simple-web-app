-- Create the role table
CREATE TABLE role
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- Add role_id column to release_artist
ALTER TABLE release_artist
    ADD COLUMN role_id INTEGER REFERENCES role (id);
