-- ---
-- Globals
-- ---

-- SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
-- SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS "matches" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;

-- ---
-- Table 'matches'
--
-- ---

DROP TABLE IF EXISTS "matches";

CREATE TABLE "matches" (
  "match_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL DEFAULT -1,
  "user_character_name" VARCHAR(100) NOT NULL,
  "user_character_gsp" INTEGER,
  "user_win" BOOLEAN,
  "opponent_character_name" VARCHAR(100) NOT NULL,
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
  "user_name" VARCHAR(100) NOT NULL,
  "email_address" VARCHAR(100) NOT NULL,
  "default_character_name" VARCHAR(100),
  "default_character_gsp" INTEGER,
  "created" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("user_id")
);


-- ---
-- Foreign Keys
-- ---

ALTER TABLE "matches" ADD FOREIGN KEY ("user_id")  REFERENCES "users" ("user_id") ON DELETE CASCADE;


-- ---
-- Table Properties
-- ---

-- ALTER TABLE "matches" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE "users" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
