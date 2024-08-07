-- Drop indexes created in the up migration
DROP INDEX IF EXISTS users_email_idx;
DROP INDEX IF EXISTS users_username_idx;

-- Drop the users table
DROP TABLE IF EXISTS "users";

-- Drop the uuid-ossp extension if it is no longer needed by other tables
DROP EXTENSION IF EXISTS "uuid-ossp";
