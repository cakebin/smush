-- First make the Tags table


DROP TABLE IF EXISTS "tags";

CREATE TABLE "tags" (
  "tag_id" SERIAL NOT NULL,
  "tag_name" VARCHAR(100) NOT NULL,
  PRIMARY KEY ("tag_id")
);


-- Populate the tags table with our existing tags
INSERT INTO "tags" ("tag_name") VALUES ('Teabagging opponent'), ('Camping opponent'), ('Homie opponent');


-- Create the match_tags table
DROP TABLE IF EXISTS "match_tags";

CREATE TABLE "match_tags" (
  "match_tag_id" SERIAL NOT NULL,
  "match_id" INTEGER NOT NULL DEFAULT -1,
  "tag_id" INTEGER NOT NULL DEFAULT -1,
  PRIMARY KEY ("match_tag_id")
);


-- Add the foreign keys to match_tags
ALTER TABLE "match_tags" ADD FOREIGN KEY ("match_id") REFERENCES "matches" ("match_id") ON DELETE CASCADE;
ALTER TABLE "match_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("tag_id") ON DELETE CASCADE;


-- Populate the match_tags for all of the matches that have "true" for each type of tag
INSERT INTO "match_tags" ("match_id", "tag_id")
SELECT "match_id", 1 AS "tag_id" FROM "matches" WHERE "opponent_teabag" = true;


INSERT INTO "match_tags" ("match_id", "tag_id")
SELECT "match_id", 2 AS "tag_id" FROM "matches" WHERE "opponent_camp" = true;


INSERT INTO "match_tags" ("match_id", "tag_id")
SELECT "match_id", 3 AS "tag_id" FROM "matches" WHERE "opponent_awesome" = true;


-- Finally, drop the old columns
ALTER TABLE "matches" DROP COLUMN "opponent_camp";
ALTER TABLE "matches" DROP COLUMN "opponent_teabag";
ALTER TABLE "matches" DROP COLUMN "opponent_awesome";
