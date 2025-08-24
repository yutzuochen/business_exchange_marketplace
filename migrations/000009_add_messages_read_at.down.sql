-- Remove read_at column from messages table
ALTER TABLE messages 
DROP COLUMN read_at;
