CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" varchar(50) UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "role" varchar(20) NOT NULL DEFAULT 'Reader',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "updated_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE UNIQUE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "users" ("username");

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
