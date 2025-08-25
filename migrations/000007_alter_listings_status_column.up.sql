-- Alter listings status column to support Chinese characters
ALTER TABLE listings 
MODIFY COLUMN status VARCHAR(50) DEFAULT '活躍';

-- Update existing records to use Chinese status
UPDATE listings SET status = '活躍' WHERE status = 'active';
UPDATE listings SET status = '不活躍' WHERE status IN ('sold', 'deleted');
