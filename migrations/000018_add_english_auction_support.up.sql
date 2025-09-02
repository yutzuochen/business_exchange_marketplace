-- 添加英式拍賣支持
-- 目標：擴展現有拍賣系統以支持英式（明標）拍賣

-- 1) 為拍賣表添加英式拍賣特定字段
ALTER TABLE auctions 
ADD COLUMN reserve_price DECIMAL(18,2) NULL COMMENT '保留價（英式拍賣用）',
ADD COLUMN min_increment DECIMAL(18,2) NOT NULL DEFAULT 10000.00 COMMENT '最小加價幅度',
ADD COLUMN buy_it_now DECIMAL(18,2) NULL COMMENT '直購價格（可選）',
ADD COLUMN current_price DECIMAL(18,2) NULL COMMENT '當前最高出價（英式拍賣用）',
ADD COLUMN highest_bidder_id BIGINT UNSIGNED NULL COMMENT '當前最高出價者',
ADD COLUMN reserve_met BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已達保留價';

-- 2) 更新出價表以支持英式拍賣的透明性
ALTER TABLE bids
ADD COLUMN max_proxy_amount DECIMAL(18,2) NULL COMMENT '代理出價上限（英式拍賣用）',
ADD COLUMN is_winning BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否為當前最高出價',
ADD COLUMN is_visible BOOLEAN NOT NULL DEFAULT TRUE COMMENT '出價是否可見（英式為true，密封為false直到結束）';

-- 3) 創建英式拍賣的實時出價視圖
CREATE INDEX idx_bids_auction_winning ON bids (auction_id, is_winning, created_at);
CREATE INDEX idx_bids_visible ON bids (auction_id, is_visible, amount DESC);

-- 4) 添加英式拍賣特定的事件類型
ALTER TABLE auction_events 
MODIFY COLUMN event_type ENUM(
  'open','bid_accepted','bid_rejected','extended','closed','notified','error',
  'reserve_met','buy_it_now','outbid','highest_bid_changed'
) NOT NULL;

-- 5) 為英式拍賣創建出價排行榜視圖（用於實時顯示）
CREATE VIEW auction_english_leaderboard AS
SELECT 
  a.auction_id,
  a.auction_type,
  a.current_price,
  a.highest_bidder_id,
  a.reserve_met,
  b.bid_id,
  b.bidder_id,
  b.amount,
  b.created_at as bid_time,
  CASE 
    WHEN a.is_anonymous THEN CONCAT('Bidder #', aba.alias_num)
    ELSE NULL 
  END as bidder_alias,
  ROW_NUMBER() OVER (PARTITION BY a.auction_id ORDER BY b.amount DESC, b.created_at ASC) as bid_rank
FROM auctions a
LEFT JOIN bids b ON a.auction_id = b.auction_id AND b.accepted = TRUE AND b.deleted_at IS NULL
LEFT JOIN auction_bidder_aliases aba ON a.auction_id = aba.auction_id AND b.bidder_id = aba.bidder_id
WHERE a.auction_type = 'english' AND b.is_visible = TRUE;

-- 6) 創建英式拍賣統計表（用於價格分布）
CREATE TABLE IF NOT EXISTS auction_english_stats (
  auction_id BIGINT UNSIGNED PRIMARY KEY,
  total_bids INT NOT NULL DEFAULT 0,
  unique_bidders INT NOT NULL DEFAULT 0,
  avg_bid_amount DECIMAL(18,2),
  median_bid_amount DECIMAL(18,2),
  bid_frequency DECIMAL(8,2) COMMENT 'bids per hour',
  last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  CONSTRAINT fk_english_stats_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 7) 注意：觸發器將在應用層實現，因為遷移工具不支持復雜的存儲過程語法