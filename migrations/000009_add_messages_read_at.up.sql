-- Add missing read_at column to messages table
ALTER TABLE messages 
ADD COLUMN read_at TIMESTAMP NULL AFTER is_read;
