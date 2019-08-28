-- First create the roles table
DROP TABLE IF EXISTS "roles";

CREATE TABLE "roles" (
  "role_id" SERIAL NOT NULL,
  "role_name" VARCHAR(100) NOT NULL,
  PRIMARY KEY ("role_id")
)


-- Create our roles
INSERT INTO "roles" (role_name) VALUES ('Admin');


-- Add a new column to users for role_id
ALTER TABLE "users" ADD COLUMN "role_id" INTEGER;


-- Update the role_id foreign key
ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id") ON DELETE CASCADE;


-- Update caebae and joraffe to be admin;
UPDATE "users" SET "role_id" = 1 WHERE "user_id" = 1;
UPDATE "users" SET "role_id" = 1 WHERE "user_id" = 2;
