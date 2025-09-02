-- 回滾英式拍賣支持

-- 1) 刪除觸發器
DROP TRIGGER IF EXISTS trg_english_auction_bid_accepted;

-- 2) 刪除英式拍賣統計表
DROP TABLE IF EXISTS auction_english_stats;

-- 3) 刪除視圖
DROP VIEW IF EXISTS auction_english_leaderboard;

-- 4) 移除出價表新增的索引
DROP INDEX IF EXISTS idx_bids_auction_winning ON bids;
DROP INDEX IF EXISTS idx_bids_visible ON bids;

-- 5) 還原事件類型枚舉
ALTER TABLE auction_events 
MODIFY COLUMN event_type ENUM(
  'open','bid_accepted','bid_rejected','extended','closed','notified','error'
) NOT NULL;

-- 6) 移除出價表新增的字段
ALTER TABLE bids
DROP COLUMN IF EXISTS max_proxy_amount,
DROP COLUMN IF EXISTS is_winning,
DROP COLUMN IF EXISTS is_visible;

-- 7) 移除拍賣表新增的字段
ALTER TABLE auctions 
DROP COLUMN IF EXISTS reserve_price,
DROP COLUMN IF EXISTS min_increment,
DROP COLUMN IF EXISTS buy_it_now,
DROP COLUMN IF EXISTS current_price,
DROP COLUMN IF EXISTS highest_bidder_id,
DROP COLUMN IF EXISTS reserve_met;