-- Remove completed_at column from transactions table
ALTER TABLE transactions 
DROP COLUMN completed_at;
