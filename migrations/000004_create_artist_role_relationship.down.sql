-- Remove the role_id column from release_artist
ALTER TABLE release_artist
    DROP COLUMN role_id;

-- Drop the role table
DROP TABLE role;
