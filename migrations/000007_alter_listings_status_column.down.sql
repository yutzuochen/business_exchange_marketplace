-- Revert listings status column back to ENUM
UPDATE listings SET status = 'active' WHERE status = '活躍';
UPDATE listings SET status = 'sold' WHERE status = '不活躍';

ALTER TABLE listings 
MODIFY COLUMN status ENUM('active', 'sold', 'deleted') DEFAULT 'active';
