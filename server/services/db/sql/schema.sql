-- ---
-- Globals
-- ---

-- SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
-- SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS "matches" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "characters" CASCADE;
DROP TABLE IF EXISTS "user_characters" CASCADE;
DROP TABLE IF EXISTS "tags" CASCADE;
DROP TABLE IF EXISTS "match_tags" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;

-- ---
-- Table 'matches'
--
-- ---

DROP TABLE IF EXISTS "matches";

CREATE TABLE "matches" (
  "match_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL DEFAULT -1,
  "user_character_id" INTEGER,
  "opponent_character_id" INTEGER NOT NULL DEFAULT -1,
  "user_character_gsp" INTEGER,
  "user_win" BOOLEAN,
  "opponent_character_gsp" INTEGER,
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
  "default_user_character_id" INTEGER,
  "user_name" VARCHAR(100) NOT NULL,
  "email_address" VARCHAR(100) NOT NULL,
  "created" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "hashed_password" VARCHAR(200) NOT NULL,
  "refresh_token" VARCHAR(200),
  "role_id" INTEGER,
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
-- Table 'user_characters'
--
-- ---

DROP TABLE IF EXISTS "user_characters";

CREATE TABLE "user_characters" (
  "user_character_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL DEFAULT -1,
  "character_id" INTEGER NOT NULL DEFAULT -1,
  "character_gsp" INTEGER,
  "alt_costume" INTEGER,
  PRIMARY KEY ("user_character_id")
);


-- ---
-- Table 'tags'
--
-- ---

DROP TABLE IF EXISTS "tags";

CREATE TABLE "tags" (
  "tag_id" SERIAL NOT NULL,
  "tag_name" VARCHAR(100) NOT NULL,
  PRIMARY KEY ("tag_id")
);


-- ---
-- Table 'match_tags'
--
-- ---

DROP TABLE IF EXISTS "match_tags";

CREATE TABLE "match_tags" (
  "match_tag_id" SERIAL NOT NULL,
  "match_id" INTEGER NOT NULL DEFAULT -1,
  "tag_id" INTEGER NOT NULL DEFAULT -1,
  PRIMARY KEY ("match_tag_id")
);


-- ---
-- Table 'roles'
--
-- ---

DROP TABLE IF EXISTS "roles";

CREATE TABLE "roles" (
  "role_id" SERIAL NOT NULL,
  "role_name" VARCHAR(100) NOT NULL,
  PRIMARY KEY ("role_id")
)

-- ---
-- Foreign Keys
-- ---

ALTER TABLE "matches" ADD FOREIGN KEY ("user_id")  REFERENCES "users" ("user_id") ON DELETE CASCADE;
ALTER TABLE "matches" ADD FOREIGN KEY ("user_character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;
ALTER TABLE "matches" ADD FOREIGN KEY ("opponent_character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;
ALTER TABLE "users" ADD FOREIGN KEY ("default_user_character_id") REFERENCES "user_characters" ("user_character_id") ON DELETE CASCADE;
ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id") ON DELETE CASCADE;
ALTER TABLE "user_characters" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;
ALTER TABLE "user_characters" ADD FOREIGN KEY ("character_id") REFERENCES "characters" ("character_id") ON DELETE CASCADE;
ALTER TABLE "match_tags" ADD FOREIGN KEY ("match_id") REFERENCES "matches" ("match_id") ON DELETE CASCADE;
ALTER TABLE "match_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("tag_id") ON DELETE CASCADE;

-- ---
-- Table Properties
-- ---

-- ALTER TABLE "matches" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE "users" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE "characters" ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
