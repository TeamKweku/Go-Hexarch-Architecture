-- Remove the index on "user_id"
DROP INDEX IF EXISTS sessions_user_id_idx;

-- Remove the foreign key constraint
ALTER TABLE "sessions" DROP CONSTRAINT IF EXISTS sessions_user_id_fkey;

-- Drop the "sessions" table
DROP TABLE IF EXISTS "sessions";
