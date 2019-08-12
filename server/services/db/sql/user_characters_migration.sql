
-- Make the user_characters table
CREATE TABLE "user_characters" (
  "user_character_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL DEFAULT -1,
  "character_id" INTEGER NOT NULL DEFAULT -1,
  "character_gsp" INTEGER,
  PRIMARY KEY ("user_character_id")
);

-- Make pre-existing entries into user_characters table for prior
--  users.default_character_id, users.default_character_gsp
-- i.e.:
INSERT INTO user_characters (user_id, character_id, character_gsp) VALUES (1, 16, 8000000);

-- Drop the users.default_character_gsp
ALTER TABLE users DROP COLUMN default_character_gsp;


-- Rename old column: users.default_character_id --> users.default_user_character_id
ALTER TABLE users RENAME COLUMN default_character_id TO default_user_character_id;

-- Add new foreign key constraint to users for user_characters
ALTER TABLE "users" ADD FOREIGN KEY ("default_user_character_id") REFERENCES "user_characters" ("user_character_id") ON DELETE CASCADE;


-- Remove the old foreign key constraint
ALTER TABLE users DROP CONSTRAINT "users_default_character_id_fkey";
