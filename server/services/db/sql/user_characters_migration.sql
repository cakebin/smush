
-- Make the user_characters table
CREATE TABLE "user_characters" (
  "user_character_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL DEFAULT -1,
  "character_id" INTEGER NOT NULL DEFAULT -1,
  "character_gsp" INTEGER,
  "alt_costume" INTEGER,
  PRIMARY KEY ("user_character_id")
);

-- Add the foreign key references in user_characters
ALTER TABLE "user_characters" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;
ALTER TABLE "user_characters" ADD FOREIGN KEY ("character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;

-- Make pre-existing entries into user_characters table for prior
--  users.default_character_id, users.default_character_gsp
-- i.e. (this will be dependant on db data):
INSERT INTO user_characters (user_id, character_id, character_gsp) VALUES (5, 39, 5200000);

-- Drop the users.default_character_gsp
ALTER TABLE users DROP COLUMN default_character_gsp;

-- Add a new column for default_user_character_id
ALTER TABLE users ADD COLUMN default_user_character_id INTEGER;

-- Add foreign key constraint to user_characters for users.default_user_character_id
ALTER TABLE "users" ADD FOREIGN KEY ("default_user_character_id") REFERENCES "user_characters" ("user_character_id") ON DELETE CASCADE;

-- Update the users.default_user_character_id to the newly inserted user_characters
-- i.e. (this will also be dependant on db data)
UPDATE users SET default_user_character_id = 2 WHERE user_id = 5;

-- Add new foreign key constraint to users for user_characters
ALTER TABLE "users" ADD FOREIGN KEY ("default_user_character_id") REFERENCES "user_characters" ("user_character_id") ON DELETE CASCADE;

-- Remove the old foreign key constraint
ALTER TABLE users DROP CONSTRAINT "users_default_character_id_fkey";

-- Drop the old column users.default_character_id
ALTER TABLE users DROP COLUMN default_character_id;
