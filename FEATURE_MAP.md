# 🏢 Business Exchange Marketplace - Comprehensive Feature Map

**Project**: 567 我來接 (567 I'll Take It)  
**Target Market**: Taiwanese business owners and entrepreneurs  
**Platform**: BizBuySell-style business exchange marketplace  
**Deployment**: Google Cloud Run + Google Cloud SQL  

---

## 📋 **Project Overview**

This is a comprehensive business exchange platform where Taiwanese entrepreneurs can:
- **Buy/Sell** existing businesses
- **Transfer** business ownership
- **List** business opportunities
- **Connect** with potential buyers/sellers
- **Manage** business transactions

---

## 🏗️ **Architecture & Technology Stack**

### **Backend (Go)**
- **Framework**: Gin HTTP web framework
- **Database**: MySQL 8 with GORM ORM
- **Cache**: Redis for performance optimization
- **Authentication**: JWT-based with bcrypt password hashing
- **API**: Dual REST API + GraphQL (gqlgen)
- **Logging**: Structured logging with Zap
- **Dependency Injection**: Wire framework
- **Containerization**: Docker with Docker Compose

### **Frontend (Next.js)**
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **State Management**: React hooks + Context
- **API Integration**: Custom API client with error handling

### **Infrastructure**
- **Hosting**: Google Cloud Run
- **Database**: Google Cloud SQL (MySQL)
- **CI/CD**: GitHub Actions
- **Monitoring**: Health checks and logging
- **Environment**: Multi-environment configuration

---

## 🔐 **Authentication & User Management**

### **User Model Features**
- **Core Fields**: ID, Email, Username, Password Hash
- **Profile**: First Name, Last Name, Phone Number
- **Security**: Role-based access (user/admin), Active status
- **Tracking**: Last login timestamp, Created/Updated dates
- **Relations**: Listings, Favorites, Messages, Transactions

### **Authentication Endpoints**
- `POST /api/v1/auth/register` - User registration with email validation
- `POST /api/v1/auth/login` - User login with JWT token generation
- **Security**: bcrypt password hashing, JWT token validation

### **User Management**
- `GET /api/v1/user/profile` - Get user profile (protected)
- `PUT /api/v1/user/profile` - Update user profile (protected)
- `PUT /api/v1/user/password` - Change password (protected)

---

## 🏪 **Business Listings System**

### **Listing Model Features**
- **Basic Info**: Title, Description, Price, Category, Condition
- **Location**: Address, Floor, Square meters
- **Business Details**: Industry, Brand story, Equipment, Decoration
- **Financial**: Annual revenue, Gross profit rate, Rent, Deposit
- **Timing**: Fastest moving date, Created/Updated timestamps
- **Status**: Active/Inactive, View count tracking
- **Relations**: Owner, Images, Favorites

### **Listing Management Endpoints**
- `GET /api/v1/listings` - List all active listings with pagination
- `GET /api/v1/listings/:id` - Get specific listing details
- `POST /api/v1/listings` - Create new listing (protected)
- `PUT /api/v1/listings/:id` - Update listing (owner only)
- `DELETE /api/v1/listings/:id` - Delete listing (owner only)

### **Advanced Search & Filtering**
- **Pagination**: Page-based with configurable limits
- **Category Filtering**: Filter by business category
- **Location Search**: Location-based filtering with LIKE queries
- **Price Range**: Min/max price filtering
- **Condition Filtering**: Filter by business condition
- **Sorting**: Default by creation date (newest first)

### **Image Management**
- `POST /api/v1/listings/:id/images` - Upload multiple images
- **Features**: Primary image designation, Order management, Alt text support
- **Storage**: Local file system with uploads directory

---

## ❤️ **Favorites & Bookmarking**

### **Favorite System**
- **Model**: User-Listing relationship with timestamps
- **Endpoints**:
  - `GET /api/v1/favorites` - List user's favorite listings
  - `POST /api/v1/favorites` - Add listing to favorites
  - `DELETE /api/v1/favorites/:id` - Remove from favorites
- **Features**: One-to-many relationship, Timestamp tracking

---

## 💬 **Messaging System**

### **Message Features**
- **Model**: Sender-Receiver relationship with content and timestamps
- **Endpoints**:
  - `GET /api/v1/messages` - List user's messages
  - `GET /api/v1/messages/:id` - Get specific message
  - `POST /api/v1/messages` - Send new message
  - `PUT /api/v1/messages/:id/read` - Mark message as read
- **Features**: Bidirectional messaging, Read status tracking

---

## 💰 **Transaction Management**

### **Transaction Model**
- **Fields**: Buyer ID, Listing ID, Amount, Status, Timestamps
- **Purpose**: Track business sales and transfers
- **Integration**: Connected to listings and users

---

## 🌐 **Web Interface (Go Templates)**

### **Server-Side Rendered Pages**
- **Homepage** (`/`): Welcome page with recent transactions and listings
- **Market Home** (`/market`): Business listing marketplace
- **Listing Detail** (`/market/listings/:id`): Individual business details
- **User Dashboard** (`/dashboard`): User management interface
- **Authentication**: Login and registration pages

### **Template Features**
- **Dynamic Content**: Database-driven listings and transactions
- **Responsive Design**: Mobile-friendly layouts
- **Search Integration**: Direct search functionality
- **Image Display**: Multi-image support for listings

---

## 🚀 **Modern Frontend (Next.js)**

### **React Components**
- **ListingCard**: Responsive business listing cards
- **Market Page**: Category filtering and listing display
- **Navigation**: Link-based routing between pages

### **API Integration**
- **Custom Client**: TypeScript-based API client
- **Error Handling**: Comprehensive error management
- **Data Fetching**: React hooks with loading states
- **Type Safety**: Full TypeScript interfaces

### **UI/UX Features**
- **Loading States**: Spinner animations during data fetch
- **Error Boundaries**: User-friendly error messages
- **Responsive Design**: Mobile-first approach
- **Modern Styling**: Tailwind CSS with hover effects

---

## 🔍 **Search & Discovery**

### **Search Features**
- **Title Search**: `/market/search?q=query` endpoint
- **Category Filtering**: Dropdown-based category selection
- **Location Filtering**: Geographic-based filtering
- **Price Range**: Min/max price sliders
- **Condition Filtering**: Business condition options

### **Discovery Features**
- **Recent Listings**: Latest business opportunities
- **Popular Categories**: Trending business types
- **Location-based**: Nearby business opportunities
- **Featured Listings**: Highlighted business opportunities

---

## 📊 **Data Management**

### **Database Features**
- **Auto-migration**: Automatic schema updates
- **Seed Data**: Sample data for testing and demonstration
- **Relationships**: Proper foreign key constraints
- **Indexing**: Performance optimization on key fields

### **Caching Strategy**
- **Redis Integration**: Session and search result caching
- **TTL Management**: Configurable cache expiration
- **Performance**: Reduced database load for frequent queries

---

## 🛡️ **Security & Middleware**

### **Security Features**
- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt with configurable cost
- **CORS Protection**: Cross-origin request security
- **Request Validation**: Input sanitization and validation

### **Middleware Stack**
- **Recovery**: Panic recovery and error handling
- **Request ID**: Unique request tracking
- **Logging**: Structured request/response logging
- **Authentication**: JWT token validation
- **CORS**: Cross-origin resource sharing

---

## 📱 **API Endpoints Summary**

### **Public Endpoints**
```
GET  /                    - Homepage
GET  /market             - Market homepage
GET  /market/search      - Search listings
GET  /market/listings/:id - Listing details
GET  /login              - Login page
GET  /register           - Registration page
GET  /healthz            - Health check
GET  /playground         - GraphQL playground
```

### **REST API v1**
```
POST   /api/v1/auth/register     - User registration
POST   /api/v1/auth/login        - User login
GET    /api/v1/listings          - List all listings
GET    /api/v1/listings/:id      - Get specific listing
GET    /api/v1/categories        - Get all categories
```

### **Protected Endpoints (JWT Required)**
```
GET    /api/v1/user/profile           - Get user profile
PUT    /api/v1/user/profile           - Update profile
PUT    /api/v1/user/password          - Change password
POST   /api/v1/listings               - Create listing
PUT    /api/v1/listings/:id           - Update listing
DELETE /api/v1/listings/:id           - Delete listing
POST   /api/v1/listings/:id/images    - Upload images
GET    /api/v1/favorites              - List favorites
POST   /api/v1/favorites              - Add favorite
DELETE /api/v1/favorites/:id          - Remove favorite
GET    /api/v1/messages               - List messages
GET    /api/v1/messages/:id           - Get message
POST   /api/v1/messages               - Send message
PUT    /api/v1/messages/:id/read      - Mark as read
```

### **GraphQL Endpoint**
```
POST /graphql - GraphQL queries and mutations
```

---

## 🚀 **Deployment & DevOps**

### **Containerization**
- **Dockerfile**: Multi-stage build for production
- **Docker Compose**: Local development environment
- **Environment Variables**: Comprehensive configuration management

### **Cloud Deployment**
- **Google Cloud Run**: Serverless backend hosting
- **Google Cloud SQL**: Managed MySQL database
- **Environment Management**: Production vs development configs

### **CI/CD Pipeline**
- **GitHub Actions**: Automated deployment
- **Build Process**: Docker image building and pushing
- **Environment Variables**: Secure secret management

---

## 📈 **Performance & Scalability**

### **Performance Features**
- **Database Connection Pooling**: Configurable connection limits
- **Redis Caching**: Search results and session caching
- **Image Optimization**: Efficient image storage and delivery
- **Pagination**: Large dataset handling

### **Scalability Features**
- **Stateless Design**: Cloud Run compatible architecture
- **Database Optimization**: Proper indexing and query optimization
- **Caching Strategy**: Redis-based performance improvement
- **Load Balancing**: Cloud Run automatic scaling

---

## 🔧 **Development & Testing**

### **Development Tools**
- **Makefile**: Build and deployment automation
- **Hot Reload**: Development server with auto-restart
- **Environment Management**: Local vs production configs
- **Database Seeding**: Sample data for development

### **Testing Features**
- **Health Checks**: Service availability monitoring
- **Error Handling**: Comprehensive error management
- **Logging**: Structured logging for debugging
- **Validation**: Input validation and sanitization

---

## 🌟 **Unique Features for Taiwanese Market**

### **Localization**
- **Traditional Chinese**: Full Chinese language support
- **Taiwanese Business Context**: Industry-specific categories
- **Local Currency**: TWD support in pricing
- **Geographic Focus**: Taiwan-specific location data

### **Business-Specific Fields**
- **Industry Categories**: Taiwanese business sectors
- **Financial Metrics**: Local business standards
- **Regulatory Compliance**: Taiwan business regulations
- **Cultural Context**: Taiwanese business practices

---

## 📝 **Maintenance & Updates**

### **How to Update This Document**
1. **Feature Changes**: Update relevant sections when adding/removing features
2. **API Updates**: Modify endpoint documentation when APIs change
3. **New Models**: Add to data model sections when creating new entities
4. **Deployment Changes**: Update infrastructure sections for new services

### **Version Control**
- **Document Version**: Update version number in header
- **Change Log**: Track major feature additions/removals
- **Review Process**: Regular review and update of feature completeness

---

## 🎯 **Future Roadmap**

### **Planned Features**
- **Real-time Notifications**: WebSocket-based updates
- **Advanced Search**: Full-text search with Elasticsearch
- **Mobile App**: React Native mobile application
- **Payment Integration**: Secure payment processing
- **Analytics Dashboard**: Business insights and metrics

### **Technical Improvements**
- **Microservices**: Service decomposition
- **Event Sourcing**: CQRS pattern implementation
- **API Versioning**: Backward compatibility management
- **Performance Monitoring**: APM integration

---

*Last Updated: December 2024*  
*Document Version: 1.0*  
*Maintained by: Development Team*
