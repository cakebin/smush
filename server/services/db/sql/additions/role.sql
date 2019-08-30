-- First create the roles table
DROP TABLE IF EXISTS "roles";

CREATE TABLE "roles" (
  "role_id" SERIAL NOT NULL,
  "role_name" VARCHAR(100) NOT NULL,
  PRIMARY KEY ("role_id")
);


-- Create our roles
INSERT INTO "roles" (role_name) VALUES ('Admin');


-- Then create the user_roles table

DROP TABLE IF EXISTS "user_roles";

CREATE TABLE "user_roles" (
  "user_role_id" SERIAL NOT NULL,
  "user_id" INTEGER NOT NULL,
  "role_id" INTEGER NOT NULL,
  PRIMARY KEY ("user_role_id")
);



-- Add foreign key constraints to "user_roles"
ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;
ALTER TABLE "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id") ON DELETE CASCADE;



-- Update caebae and joraffe to be admin;
INSERT INTO "user_roles" ("user_id", "role_id") VALUES (1, 1), (2, 1);
