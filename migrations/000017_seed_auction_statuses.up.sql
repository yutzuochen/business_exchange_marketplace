-- Seed auction status reference data

INSERT IGNORE INTO auction_status_ref (status_code, is_open, description) VALUES
('draft', false, 'Draft auction, not yet published'),
('active', true, 'Active auction accepting bids'),
('extended', true, 'Auction extended due to soft close'),
('closed', false, 'Auction closed, results available'),
('cancelled', false, 'Auction cancelled by seller');