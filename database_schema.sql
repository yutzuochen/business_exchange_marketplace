-- Business Marketplace Database Schema
-- MySQL 8.0+ compatible

-- Enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;

-- Users table for authentication and profile management
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) UNIQUE NOT NULL DEFAULT (UUID()),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    company_name VARCHAR(255),
    
    -- Profile information
    bio TEXT,
    profile_image_url VARCHAR(500),
    website_url VARCHAR(500),
    linkedin_url VARCHAR(500),
    
    -- Location
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'US',
    zip_code VARCHAR(20),
    
    -- Account status
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    is_premium BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP NULL,
    
    INDEX idx_email (email),
    INDEX idx_location (city, state),
    INDEX idx_created_at (created_at),
    INDEX idx_active (is_active)
);

-- Business categories for organizing listings
CREATE TABLE categories (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_id INT NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_slug (slug),
    INDEX idx_parent (parent_id),
    INDEX idx_active (is_active)
);

-- Business listings - main entity
CREATE TABLE business_listings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) UNIQUE NOT NULL DEFAULT (UUID()),
    user_id BIGINT NOT NULL,
    category_id INT NOT NULL,
    
    -- Basic information
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    short_description VARCHAR(500),
    
    -- Business details
    business_name VARCHAR(255),
    industry VARCHAR(100),
    established_year YEAR,
    number_of_employees INT,
    
    -- Financial information
    asking_price DECIMAL(15,2),
    annual_revenue DECIMAL(15,2),
    annual_profit DECIMAL(15,2),
    cash_flow DECIMAL(15,2),
    
    -- Location
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) DEFAULT 'US',
    zip_code VARCHAR(20),
    full_address TEXT,
    
    -- Listing details
    listing_type ENUM('business_for_sale', 'franchise', 'business_opportunity') DEFAULT 'business_for_sale',
    status ENUM('draft', 'active', 'pending', 'sold', 'expired', 'withdrawn') DEFAULT 'draft',
    
    -- Features and amenities (JSON array)
    features JSON,
    
    -- Images
    main_image_url VARCHAR(500),
    images JSON, -- Array of image URLs
    
    -- SEO and marketing
    meta_title VARCHAR(255),
    meta_description VARCHAR(500),
    
    -- Engagement metrics
    view_count INT DEFAULT 0,
    inquiry_count INT DEFAULT 0,
    favorite_count INT DEFAULT 0,
    
    -- Listing management
    is_featured BOOLEAN DEFAULT FALSE,
    is_urgent BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    
    -- Search and filter indexes
    INDEX idx_location (city, state),
    INDEX idx_price (asking_price),
    INDEX idx_category (category_id),
    INDEX idx_status (status),
    INDEX idx_type (listing_type),
    INDEX idx_featured (is_featured),
    INDEX idx_created (created_at),
    INDEX idx_published (published_at),
    INDEX idx_user (user_id),
    INDEX idx_slug (slug),
    
    -- Composite indexes for common searches
    INDEX idx_location_category (city, state, category_id),
    INDEX idx_price_location (asking_price, city, state),
    INDEX idx_status_featured (status, is_featured),
    
    -- Full-text search index
    FULLTEXT INDEX idx_search (title, description, business_name, industry)
);

-- User inquiries about listings
CREATE TABLE inquiries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    listing_id BIGINT NOT NULL,
    inquirer_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,
    
    -- Inquiry details
    subject VARCHAR(255),
    message TEXT NOT NULL,
    phone VARCHAR(20),
    
    -- Inquiry type
    inquiry_type ENUM('general', 'financing', 'viewing', 'offer') DEFAULT 'general',
    
    -- Status tracking
    status ENUM('new', 'read', 'replied', 'closed') DEFAULT 'new',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    read_at TIMESTAMP NULL,
    replied_at TIMESTAMP NULL,
    
    FOREIGN KEY (listing_id) REFERENCES business_listings(id) ON DELETE CASCADE,
    FOREIGN KEY (inquirer_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_listing (listing_id),
    INDEX idx_inquirer (inquirer_id),
    INDEX idx_seller (seller_id),
    INDEX idx_status (status),
    INDEX idx_created (created_at)
);

-- User favorites/watchlist
CREATE TABLE favorites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    listing_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (listing_id) REFERENCES business_listings(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_favorite (user_id, listing_id),
    INDEX idx_user (user_id),
    INDEX idx_listing (listing_id)
);

-- Search queries for analytics
CREATE TABLE search_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NULL,
    query_text VARCHAR(500),
    filters JSON, -- Location, price range, category filters
    results_count INT DEFAULT 0,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user (user_id),
    INDEX idx_created (created_at),
    INDEX idx_query (query_text(100))
);

-- Email verification tokens
CREATE TABLE email_verifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_token (token),
    INDEX idx_user (user_id),
    INDEX idx_expires (expires_at)
);

-- Password reset tokens
CREATE TABLE password_resets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_token (token),
    INDEX idx_user (user_id),
    INDEX idx_expires (expires_at)
);

-- Insert sample categories
INSERT INTO categories (name, slug, description, sort_order) VALUES
('Restaurants & Food Service', 'restaurants-food-service', 'Restaurants, cafes, catering, and food service businesses', 1),
('Retail & E-commerce', 'retail-ecommerce', 'Retail stores, online businesses, and e-commerce platforms', 2),
('Healthcare & Medical', 'healthcare-medical', 'Medical practices, clinics, and healthcare services', 3),
('Technology & Software', 'technology-software', 'Tech companies, software businesses, and IT services', 4),
('Manufacturing & Industrial', 'manufacturing-industrial', 'Manufacturing, production, and industrial businesses', 5),
('Professional Services', 'professional-services', 'Consulting, legal, accounting, and other professional services', 6),
('Beauty & Wellness', 'beauty-wellness', 'Salons, spas, fitness centers, and wellness businesses', 7),
('Automotive', 'automotive', 'Auto repair, dealerships, and automotive services', 8),
('Real Estate', 'real-estate', 'Real estate agencies, property management, and related services', 9),
('Education & Training', 'education-training', 'Schools, training centers, and educational services', 10),
('Entertainment & Media', 'entertainment-media', 'Entertainment venues, media companies, and creative services', 11),
('Franchises', 'franchises', 'Franchise opportunities across various industries', 12);

-- Create indexes for better performance
CREATE INDEX idx_listings_search ON business_listings (status, is_featured, asking_price, created_at);
CREATE INDEX idx_listings_location_price ON business_listings (city, state, asking_price) WHERE status = 'active';
CREATE INDEX idx_users_location ON users (city, state) WHERE is_active = TRUE; 