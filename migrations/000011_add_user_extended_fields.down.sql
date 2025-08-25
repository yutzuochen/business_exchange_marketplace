-- Remove extended user fields
DROP INDEX IF EXISTS idx_users_two_factor_enabled ON users;
DROP INDEX IF EXISTS idx_users_email_verified_at ON users;

ALTER TABLE users 
DROP COLUMN IF EXISTS marketing_emails,
DROP COLUMN IF EXISTS email_notifications,
DROP COLUMN IF EXISTS contact_phone,
DROP COLUMN IF EXISTS tax_id,
DROP COLUMN IF EXISTS company_name,
DROP COLUMN IF EXISTS two_factor_secret,
DROP COLUMN IF EXISTS two_factor_enabled,
DROP COLUMN IF EXISTS email_verification_token,
DROP COLUMN IF EXISTS email_verified_at;
