-- Rollback auction core tables in reverse dependency order
-- This allows clean rollback if needed

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS auction_stream_offsets;
DROP TABLE IF EXISTS auction_notification_log;
DROP TABLE IF EXISTS user_blacklist;
DROP TABLE IF EXISTS auction_bid_histograms;
DROP TABLE IF EXISTS auction_bidder_aliases;
DROP TABLE IF EXISTS auction_events;
DROP TABLE IF EXISTS auction_status_history;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS auctions;
DROP TABLE IF EXISTS auction_status_ref;