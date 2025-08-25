-- Add extended user fields for email verification, 2FA, and seller info
ALTER TABLE users 
ADD COLUMN email_verified_at TIMESTAMP NULL,
ADD COLUMN email_verification_token VARCHAR(255) DEFAULT '',
ADD COLUMN two_factor_enabled BOOLEAN DEFAULT FALSE,
ADD COLUMN two_factor_secret VARCHAR(255) DEFAULT '',
ADD COLUMN company_name VARCHAR(255) DEFAULT '',
ADD COLUMN tax_id VARCHAR(20) DEFAULT '',
ADD COLUMN contact_phone VARCHAR(20) DEFAULT '',
ADD COLUMN email_notifications BOOLEAN DEFAULT TRUE,
ADD COLUMN marketing_emails BOOLEAN DEFAULT FALSE;
