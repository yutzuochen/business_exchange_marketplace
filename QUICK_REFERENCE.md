# 🚀 Quick Reference Guide - Business Exchange Marketplace

## 📍 **Project Quick Facts**
- **Backend**: Go + Gin + MySQL + Redis
- **Frontend**: Next.js + TypeScript + Tailwind CSS
- **Deployment**: Google Cloud Run + Cloud SQL
- **Brand**: 567 我來接 (567 I'll Take It)
- **Branch**: feat/nextJS

## 🔑 **Key Commands**

### **Backend Development**
```bash
# Run locally
make run
go run ./cmd/server

# Build
make build

# Docker
make docker-up
make docker-down

# Generate GraphQL code
make gqlgen
```

### **Frontend Development**
```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## 🌐 **Main URLs**

### **Local Development**
- **Backend**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Database Admin**: http://localhost:8081 (Adminer)

### **Production**
- **Backend**: https://business-exchange-backend-430730011391.us-central1.run.app
- **Frontend**: [Your frontend URL]

## 📱 **API Endpoints**

### **Most Used**
- `GET /api/v1/listings` - Get all listings
- `GET /api/v1/listings/:id` - Get specific listing
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/listings` - Create listing (auth required)

### **Authentication**
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### **Health Check**
- `GET /healthz` - Service health status

## 🗄️ **Database Models**

### **Core Entities**
- **User**: Authentication, profile, role management
- **Listing**: Business listings with detailed attributes
- **Image**: Multi-image support for listings
- **Favorite**: User bookmarking system
- **Message**: Internal messaging between users
- **Transaction**: Business sale/transfer tracking

### **Key Fields**
- **Listing**: title, price, category, location, industry, annual_revenue
- **User**: email, username, role, is_active
- **Image**: filename, url, is_primary, order

## 🔧 **Configuration**

### **Environment Variables**
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=business_exchange

# Redis
REDIS_ADDR=localhost:6379

# JWT
JWT_SECRET=your-secret-key
JWT_ISSUER=business-exchange

# App
APP_PORT=8080
APP_ENV=development
```

## 📁 **Key Directories**

### **Backend**
- `cmd/server/` - Main application entry point
- `internal/` - Core business logic
- `templates/` - HTML templates
- `static/` - Static assets
- `uploads/` - User uploaded images

### **Frontend**
- `src/app/` - Next.js app router pages
- `src/components/` - React components
- `src/types/` - TypeScript interfaces
- `src/lib/` - Utility functions and API client

## 🚨 **Common Issues & Solutions**

### **Database Connection**
```bash
# Test connection
./test-db-connection.go

# Check Cloud SQL
./test-cloud-sql.sh
```

### **Environment Variables**
```bash
# Fix environment variables
./fix-env-vars.sh
./fix-env-vars-simple.sh
```

### **Deployment Issues**
```bash
# Fix Cloud Run
./fix-cloud-run.sh

# Deploy to cloud
./deploy-to-cloud.sh
```

## 📊 **Data Examples**

### **Sample Listing**
```json
{
  "title": "台北市咖啡廳轉讓",
  "description": "位於信義區的精美咖啡廳",
  "price": 2500000,
  "category": "餐飲業",
  "location": "台北市信義區",
  "industry": "餐飲服務",
  "annual_revenue": 8000000,
  "square_meters": 45.5
}
```

### **Sample User**
```json
{
  "email": "user@example.com",
  "username": "businessowner",
  "first_name": "王",
  "last_name": "小明",
  "role": "user",
  "is_active": true
}
```

## 🔍 **Search & Filtering**

### **Listing Search Parameters**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)
- `category` - Business category
- `location` - Location filter
- `min_price` - Minimum price
- `max_price` - Maximum price
- `condition` - Business condition

### **Example Search**
```
GET /api/v1/listings?category=餐飲業&location=台北市&min_price=1000000&max_price=5000000
```

## 🎨 **Frontend Components**

### **Main Components**
- `ListingCard` - Business listing display card
- `MarketPage` - Main marketplace page
- `ListingDetail` - Individual listing view

### **Styling**
- **Framework**: Tailwind CSS v4
- **Theme**: Blue/Orange color scheme
- **Responsive**: Mobile-first design

## 📝 **Development Notes**

### **Code Style**
- **Backend**: Go standard formatting
- **Frontend**: ESLint + Prettier
- **Database**: GORM conventions

### **Testing**
- **Health Checks**: `/healthz` endpoint
- **Database**: Auto-migration + seed data
- **API**: REST + GraphQL support

## 🚀 **Deployment Checklist**

### **Before Deploying**
- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] Docker image builds successfully
- [ ] Health checks pass
- [ ] API endpoints respond correctly

### **After Deployment**
- [ ] Verify health check endpoint
- [ ] Test authentication flow
- [ ] Verify database connections
- [ ] Check image uploads work
- [ ] Test search functionality

## 🔄 **Branch-Specific Notes**

### **feat/nextJS Branch**
- **Focus**: Next.js frontend integration
- **Status**: Development in progress
- **Key Changes**: Modern React components, TypeScript interfaces
- **Integration**: API client for backend communication

---

*This is a living document - update as needed during development*  
*Branch: feat/nextJS*
