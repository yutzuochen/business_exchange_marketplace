# 🚀 Members Service - Business Exchange Marketplace

## 📋 Overview

The Members Service is a comprehensive authentication and user management system integrated into the Business Exchange Marketplace backend. It provides secure, session-based authentication with advanced security features, email verification, and lead management capabilities.

## ✨ Key Features

### 🔐 **Authentication & Security**
- **Session-based authentication** (no JWT for browser)
- **Email verification** before account activation
- **Password reset** with secure tokens
- **Rate limiting** for login, signup, and contact forms
- **Account lockout** after failed login attempts
- **Anti-bot protection** with honeypot fields and timing checks

### 👥 **User Management**
- **Role-based access control** (user, seller, admin)
- **Seller-specific fields** (company name, tax ID, contact phone)
- **Profile management** with notification preferences
- **Active sessions tracking** with revoke capability

### 📧 **Email Services**
- **Email verification** for new accounts
- **Password reset emails** with secure tokens
- **Lead notifications** for sellers
- **SendGrid integration** (production) / Logging (development)

### 💼 **Lead Management**
- **Contact seller forms** with anti-spam measures
- **Lead tracking** and notification system
- **Spam detection** with keyword filtering
- **Rate limiting** per user and seller

## 🏗️ Architecture

### **Core Components**

```
internal/
├── auth/
│   ├── session.go      # Session management with Redis
│   ├── email.go        # Email service (SendGrid/logging)
│   └── jwt.go          # JWT for API tokens (legacy)
├── models/
│   └── user.go         # Enhanced user model with new fields
├── middleware/
│   ├── session.go      # Session-based authentication
│   ├── rate_limit.go   # Rate limiting middleware
│   └── jwt.go          # JWT middleware (legacy)
├── handlers/
│   ├── members_auth.go # Enhanced auth handlers
│   └── leads.go        # Lead management handlers
└── config/
    └── config.go       # Enhanced configuration
```

### **Database Models**

- **User**: Enhanced with email verification, 2FA, seller fields
- **UserSession**: Server-side session management
- **Lead**: Contact form submissions and lead tracking
- **PasswordResetToken**: Secure password reset tokens
- **AuditLog**: Security event logging

## 🚀 Getting Started

### **1. Environment Variables**

Add these to your `.env` file:

```bash
# Members Service Configuration
SENDGRID_API_KEY=your_sendgrid_api_key
SENDGRID_FROM_EMAIL=noreply@yourdomain.com
SENDGRID_FROM_NAME=Business Exchange

# Session Management
SESSION_SECRET=your-super-secret-session-key
SESSION_TTL_MINUTES=1440
SESSION_COOKIE_DOMAIN=.yourdomain.com
SESSION_COOKIE_SECURE=true
SESSION_COOKIE_HTTP_ONLY=true
SESSION_COOKIE_SAME_SITE=Lax

# Rate Limiting
RATE_LIMIT_LOGIN_PER_MINUTE=5
RATE_LIMIT_SIGNUP_PER_HOUR=3
RATE_LIMIT_FORGOT_PASSWORD_PER_HOUR=3
RATE_LIMIT_CONTACT_SELLER_PER_HOUR=10

# Security
PASSWORD_MIN_LENGTH=8
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION_MINUTES=30

# File Upload Limits
MAX_FILE_SIZE_MB=5
MAX_TOTAL_SIZE_MB=25
MAX_FILES_PER_REQUEST=10
MAX_AVATAR_SIZE_MB=1
GLOBAL_BODY_LIMIT_MB=30
```

### **2. Database Migration**

The new models will be automatically migrated when you start the application. Ensure your database supports the new fields.

### **3. Redis Setup**

Redis is required for session storage and rate limiting:

```bash
# Local development
docker run -d -p 6379:6379 redis:alpine

# Or use existing Redis instance
REDIS_ADDR=localhost:6379
```

## 📡 API Endpoints

### **Authentication**

```http
POST /api/v1/auth/signup
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/auth/verify-email
POST /api/v1/auth/forgot-password
POST /api/v1/auth/reset-password
```

### **User Management**

```http
GET  /api/v1/users/profile
PUT  /api/v1/users/profile
GET  /api/v1/users/sessions
DELETE /api/v1/users/sessions/:session_id
```

### **Lead Management**

```http
POST /api/v1/leads/contact-seller
GET  /api/v1/leads
PUT  /api/v1/leads/:id/read
```

### **Admin Endpoints**

```http
GET  /api/v1/admin/users
GET  /api/v1/admin/users/:id
PUT  /api/v1/admin/users/:id/status
PUT  /api/v1/admin/users/:id/role
GET  /api/v1/admin/leads
```

## 🔒 Security Features

### **Anti-Bot Protection**

1. **Honeypot Fields**: Hidden form fields that bots fill out
2. **Timing Checks**: Reject forms submitted too quickly (< 800ms)
3. **Cloudflare Turnstile**: CAPTCHA alternative (production)

### **Rate Limiting**

- **Login**: 5 attempts per minute per IP
- **Signup**: 3 attempts per hour per IP
- **Password Reset**: 3 requests per hour per email
- **Contact Seller**: 10 requests per hour per user-seller pair

### **Session Security**

- **HttpOnly cookies**: Prevents XSS attacks
- **Secure cookies**: HTTPS only in production
- **SameSite=Lax**: CSRF protection
- **Automatic expiration**: Configurable TTL
- **Session revocation**: Users can revoke individual sessions

## 📧 Email Templates

### **Verification Email**
- Welcome message with verification link
- 24-hour expiration
- Professional HTML and text versions

### **Password Reset Email**
- Secure reset link
- 30-minute expiration
- Clear instructions

### **Lead Notification**
- Lead details for sellers
- Contact information
- Professional formatting

## 🛡️ Frontend Integration

### **Session Cookies**

The service automatically sets secure session cookies. Frontend should:

```javascript
// Always include credentials
fetch('/api/v1/users/profile', {
  credentials: 'include'
})

// Next.js with credentials
const response = await fetch('/api/v1/users/profile', {
  credentials: 'include'
})
```

### **Anti-Bot Fields**

Include these in your forms:

```html
<!-- Hidden honeypot field -->
<input type="text" name="website" style="display: none;" tabindex="-1">

<!-- Form render time -->
<input type="hidden" name="form_time" value="<%= Date.now() %>">
```

### **Cloudflare Turnstile**

For production, add Turnstile to your forms:

```html
<div class="cf-turnstile" data-sitekey="your-site-key"></div>
```

## 🔧 Configuration Options

### **Development vs Production**

- **Development**: Emails are logged to console, cookies less restrictive
- **Production**: SendGrid emails, strict cookie security, Turnstile verification

### **Customization**

All security parameters are configurable via environment variables:

- Rate limiting thresholds
- Session timeouts
- File upload limits
- Password requirements

## 📊 Monitoring & Logging

### **Audit Events**

The system logs security events:

- Login attempts (success/failure)
- Password changes
- Role changes
- Session creation/revocation

### **Rate Limit Monitoring**

Redis stores rate limit counters for monitoring and debugging.

## 🚨 Troubleshooting

### **Common Issues**

1. **Sessions not persisting**: Check Redis connection and cookie settings
2. **Rate limiting too strict**: Adjust environment variables
3. **Emails not sending**: Verify SendGrid configuration
4. **Cookie issues**: Check domain and secure settings

### **Debug Mode**

Set `APP_ENV=development` for detailed logging and relaxed security.

## 🔮 Future Enhancements

- **2FA Support**: TOTP-based two-factor authentication
- **Social Login**: OAuth integration
- **Advanced Spam Detection**: Machine learning-based filtering
- **Audit Dashboard**: Admin interface for security monitoring
- **Webhook Support**: Real-time notifications

## 📚 Additional Resources

- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [SendGrid API Documentation](https://sendgrid.com/docs/api-reference/)

---

*This service provides enterprise-grade security while maintaining developer-friendly configuration and comprehensive documentation.*
