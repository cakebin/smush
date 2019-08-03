-- ---
-- Globals
-- ---

-- SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
-- SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS "matches" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "characters" CASCADE;

-- ---
-- Table 'matches'
--
-- ---

DROP TABLE IF EXISTS "matches";

CREATE TABLE "matches" (
  "match_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL DEFAULT -1,
  "user_character_id" INTEGER NOT NULL DEFAULT -1,
  "opponent_character_id" INTEGER NOT NULL DEFAULT -1,
  "user_character_gsp" INTEGER,
  "user_win" BOOLEAN,
  "opponent_character_gsp" INTEGER,
  "opponent_teabag" BOOLEAN,
  "opponent_camp" BOOLEAN,
  "opponent_awesome" BOOLEAN,
  "created" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("match_id")
);


-- ---
-- Table 'users'
--
-- ---

DROP TABLE IF EXISTS "users";

CREATE TABLE "users" (
  "user_id" SERIAL NOT NULL,
  "default_character_id" INTEGER,
  "user_name" VARCHAR(100) NOT NULL,
  "email_address" VARCHAR(100) NOT NULL,
  "default_character_gsp" INTEGER,
  "created" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "hashed_password" VARCHAR(200) NOT NULL,
  "refresh_token" VARCHAR(200),
  PRIMARY KEY ("user_id")
);


-- ---
-- Table 'users'
--
-- ---

DROP TABLE IF EXISTS "characters";

CREATE TABLE "characters" (
  "character_id" SERIAL NOT NULL,
  "character_name" VARCHAR(100) NOT NULL,
  "character_stock_img" VARCHAR(100),
  "character_img" VARCHAR(100),
  "character_archetype" VARCHAR(100),
  PRIMARY KEY ("character_id")
);


-- ---
-- Foreign Keys
-- ---

ALTER TABLE "matches" ADD FOREIGN KEY ("user_id")  REFERENCES "users" ("user_id") ON DELETE CASCADE;
ALTER TABLE "matches" ADD FOREIGN KEY ("user_character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;
ALTER TABLE "matches" ADD FOREIGN KEY ("opponent_character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;
ALTER TABLE "users" ADD FOREIGN KEY ("default_character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;


-- ---
-- Table Properties
-- ---

-- ALTER TABLE "matches" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE "users" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE "characters" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
