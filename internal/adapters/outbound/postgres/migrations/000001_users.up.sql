CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "etag" varchar(250) NOT NULL,
  "username" varchar(50) NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "role" varchar(20) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("username");
