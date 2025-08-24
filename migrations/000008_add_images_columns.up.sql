-- Add missing columns to images table
ALTER TABLE images 
ADD COLUMN filename VARCHAR(255) NOT NULL AFTER listing_id,
ADD COLUMN `order` INT DEFAULT 0 AFTER alt_text;
