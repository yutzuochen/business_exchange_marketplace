-- Create leads table for contact form submissions
CREATE TABLE leads (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    receiver_id BIGINT NOT NULL,
    listing_id BIGINT NULL,
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    contact_phone VARCHAR(20),
    is_read BOOLEAN DEFAULT FALSE,
    is_spam BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_leads_sender_id (sender_id),
    INDEX idx_leads_receiver_id (receiver_id),
    INDEX idx_leads_listing_id (listing_id),
    INDEX idx_leads_is_read (is_read),
    INDEX idx_leads_is_spam (is_spam),
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE SET NULL
);
