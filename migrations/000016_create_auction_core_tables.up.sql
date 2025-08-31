-- 拍賣系統核心表（不包含 audit_logs，因為已存在）
-- 目標：MySQL 8.x / InnoDB / utf8mb4

-- 1) 參考狀態表
CREATE TABLE IF NOT EXISTS auction_status_ref (
  status_code VARCHAR(16) PRIMARY KEY,
  is_open BOOLEAN NOT NULL,
  description VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2) 拍賣主表（密封投標為預設）
CREATE TABLE IF NOT EXISTS auctions (
  auction_id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  listing_id BIGINT UNSIGNED NOT NULL,
  seller_id  BIGINT UNSIGNED NOT NULL,
  auction_type ENUM('sealed','english','dutch') NOT NULL DEFAULT 'sealed',
  status_code VARCHAR(16) NOT NULL,
  allowed_min_bid DECIMAL(18,2) NOT NULL,
  allowed_max_bid DECIMAL(18,2) NOT NULL,
  soft_close_trigger_sec INT NOT NULL DEFAULT 180,  -- 結束前 3 分鐘觸發
  soft_close_extend_sec  INT NOT NULL DEFAULT 60,   -- 延長 1 分鐘
  start_at   DATETIME NOT NULL,
  end_at     DATETIME NOT NULL,
  extended_until DATETIME NULL,
  extension_count INT NOT NULL DEFAULT 0,
  is_anonymous BOOLEAN NOT NULL DEFAULT TRUE,
  view_count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT chk_bid_range CHECK (allowed_min_bid >= 0 AND allowed_max_bid > allowed_min_bid),
  CONSTRAINT chk_duration  CHECK (TIMESTAMPDIFF(DAY, start_at, end_at) BETWEEN 1 AND 61),

  CONSTRAINT fk_auction_status
    FOREIGN KEY (status_code) REFERENCES auction_status_ref(status_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_auctions_status_end ON auctions (status_code, end_at);
CREATE INDEX idx_auctions_listing     ON auctions (listing_id);
CREATE INDEX idx_auctions_seller      ON auctions (seller_id);

-- 3) 出價表（盲標、軟刪除、結束時計名）
CREATE TABLE IF NOT EXISTS bids (
  bid_id     BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  auction_id BIGINT UNSIGNED NOT NULL,
  bidder_id  BIGINT UNSIGNED NOT NULL,
  amount     DECIMAL(18,2) NOT NULL,
  client_seq BIGINT NOT NULL,
  source_ip_hash VARBINARY(32) NULL,
  user_agent_hash VARBINARY(32) NULL,
  accepted   BOOLEAN NOT NULL DEFAULT TRUE,
  reject_reason VARCHAR(64) NULL,            -- 'out_of_range'/'too_late'/'blacklisted'/...
  final_rank INT NULL,                        -- 拍賣結束後寫入
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  deleted_by BIGINT UNSIGNED NULL,

  CONSTRAINT fk_bids_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id),
  CONSTRAINT chk_bid_amount  CHECK (amount >= 0),

  UNIQUE KEY uk_idem (auction_id, bidder_id, client_seq),
  INDEX idx_bids_user (bidder_id, auction_id, created_at),
  INDEX idx_bids_auction_time (auction_id, created_at),
  INDEX idx_bids_auction_amount (auction_id, amount DESC),
  INDEX idx_bids_auction_rank (auction_id, final_rank)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4) 拍賣狀態歷史
CREATE TABLE IF NOT EXISTS auction_status_history (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  auction_id BIGINT UNSIGNED NOT NULL,
  from_status VARCHAR(16) NOT NULL,
  to_status   VARCHAR(16) NOT NULL,
  reason      VARCHAR(255) NULL,
  operator_id BIGINT UNSIGNED NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_hist_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_hist_auction_time ON auction_status_history (auction_id, created_at);

-- 5) 拍賣事件（WS 對帳、斷線恢復、審計）
CREATE TABLE IF NOT EXISTS auction_events (
  event_id   BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  auction_id BIGINT UNSIGNED NOT NULL,
  event_type ENUM('open','bid_accepted','bid_rejected','extended','closed','notified','error') NOT NULL,
  actor_user_id BIGINT UNSIGNED NULL,
  payload JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_events_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_events_auction_time ON auction_events (auction_id, created_at);

-- 6) 匿名別名（Bidder #23）
CREATE TABLE IF NOT EXISTS auction_bidder_aliases (
  auction_id BIGINT UNSIGNED NOT NULL,
  bidder_id  BIGINT UNSIGNED NOT NULL,
  alias_num  INT NOT NULL,
  alias_label VARCHAR(32) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (auction_id, bidder_id),
  UNIQUE KEY uk_alias_label (auction_id, alias_label),
  CONSTRAINT fk_alias_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 7) 出價分佈快照（背景任務每 5 分鐘寫一次）
CREATE TABLE IF NOT EXISTS auction_bid_histograms (
  auction_id BIGINT UNSIGNED NOT NULL,
  bucket_low  DECIMAL(18,2) NOT NULL,
  bucket_high DECIMAL(18,2) NOT NULL,
  bid_count   INT NOT NULL,
  computed_at DATETIME NOT NULL,

  PRIMARY KEY (auction_id, bucket_low, bucket_high, computed_at),
  CONSTRAINT fk_histogram_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 8) 黑名單（全站）
CREATE TABLE IF NOT EXISTS user_blacklist (
  user_id BIGINT UNSIGNED PRIMARY KEY,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  reason VARCHAR(255) NULL,
  staff_id BIGINT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 9) 通知紀錄（得標者/前7名/參與者/賣家）
CREATE TABLE IF NOT EXISTS auction_notification_log (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  auction_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  kind ENUM('winner','seller_result','top7','participant_end') NOT NULL,
  channel ENUM('email','sms','line','webpush') NOT NULL,
  status ENUM('queued','sent','failed') NOT NULL DEFAULT 'queued',
  meta JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE KEY uk_once (auction_id, user_id, kind),
  INDEX idx_notif_auction (auction_id, created_at),
  CONSTRAINT fk_notif_auction FOREIGN KEY (auction_id) REFERENCES auctions(auction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 10) WS 斷線恢復（每用戶最後讀到的事件）
CREATE TABLE IF NOT EXISTS auction_stream_offsets (
  auction_id BIGINT UNSIGNED NOT NULL,
  user_id    BIGINT UNSIGNED NOT NULL,
  last_event_id BIGINT UNSIGNED NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (auction_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;