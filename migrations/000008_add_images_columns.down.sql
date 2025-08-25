-- Remove added columns from images table
ALTER TABLE images 
DROP COLUMN filename,
DROP COLUMN `order`;
