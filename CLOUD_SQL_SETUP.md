# 🚀 Google Cloud SQL 設置指南

## 📋 專案資訊
- **專案名稱**: BusinessExchange
- **專案 ID**: businessexchange-468413
- **資料庫類型**: Cloud SQL (MySQL)

## 🔧 環境變數配置

### 1. 創建生產環境 `.env` 文件

```bash
# 複製範例文件
cp .env.example .env.production
```

### 2. 更新 `.env.production` 文件

```env
# Production Environment Configuration
APP_ENV=production
APP_PORT=8080
APP_NAME=BusinessExchange

# Cloud SQL Configuration
DB_HOST=YOUR_CLOUD_SQL_IP_ADDRESS
DB_PORT=3306
DB_USER=YOUR_DATABASE_USERNAME
DB_PASSWORD=YOUR_DATABASE_PASSWORD
DB_NAME=business_exchange
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_LOC=Local
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100

# Redis Configuration (Cloud Memorystore)
REDIS_ADDR=YOUR_MEMORYSTORE_IP:6379
REDIS_DB=0
REDIS_PASSWORD=
REDIS_POOL_SIZE=10

# JWT Configuration
JWT_SECRET=YOUR_SUPER_SECURE_JWT_SECRET_KEY_HERE
JWT_ISSUER=businessexchange-468413
JWT_EXPIRE_MINUTES=60
JWT_REFRESH_EXPIRE_DAYS=7

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://www.your-domain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With

# File Upload Configuration
MAX_FILE_SIZE=10485760
UPLOAD_DIR=/tmp/uploads
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/gif,image/webp

# Pagination
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100
```

## 🌐 Cloud SQL 連接資訊

### 獲取連接資訊

1. **登入 Google Cloud Console**
   - 前往: https://console.cloud.google.com/
   - 選擇專案: `businessexchange-468413`

2. **找到 Cloud SQL 實例**
   - 前往: SQL > 實例
   - 點擊您的 MySQL 實例

3. **獲取連接資訊**
   - **公共 IP 地址**: 在 "概覽" 頁面查看
   - **連接名稱**: 在 "概覽" 頁面查看
   - **資料庫版本**: MySQL 8.0

### 資料庫用戶設置

1. **創建資料庫用戶**
   ```sql
   CREATE USER 'app'@'%' IDENTIFIED BY 'your_secure_password';
   GRANT ALL PRIVILEGES ON business_exchange.* TO 'app'@'%';
   FLUSH PRIVILEGES;
   ```

2. **創建資料庫**
   ```sql
   CREATE DATABASE business_exchange CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

## 🔐 安全配置

### 1. 授權網路

在 Cloud SQL 實例的 "連線" 頁面：
- 添加您的應用程式 IP 地址
- 或者添加 `0.0.0.0/0` (僅用於測試，生產環境不建議)

### 2. SSL 連接

Cloud SQL 預設啟用 SSL，確保在應用程式中使用 SSL 連接。

## 🚀 部署配置

### 1. Cloud Run 部署

```bash
# 構建 Docker 映像
docker build -t gcr.io/businessexchange-468413/business-exchange .

# 推送到 Google Container Registry
docker push gcr.io/businessexchange-468413/business-exchange

# 部署到 Cloud Run
gcloud run deploy business-exchange \
  --image gcr.io/businessexchange-468413/business-exchange \
  --platform managed \
  --region asia-east1 \
  --allow-unauthenticated \
  --set-env-vars "APP_ENV=production"
```

### 2. 環境變數設置

在 Cloud Run 服務中設置環境變數：
- `DB_HOST`: 您的 Cloud SQL IP
- `DB_USER`: 資料庫用戶名
- `DB_PASSWORD`: 資料庫密碼
- `JWT_SECRET`: 安全的 JWT 密鑰

## 📊 監控和日誌

### 1. Cloud Logging
- 應用程式日誌會自動發送到 Cloud Logging
- 可以在 Google Cloud Console 中查看

### 2. Cloud Monitoring
- 設置 Cloud SQL 監控
- 監控資料庫性能和連接數

## 🔍 測試連接

### 1. 本地測試
```bash
# 使用 Cloud SQL Proxy 進行本地測試
cloud_sql_proxy -instances=businessexchange-468413:asia-east1:your-instance-name=tcp:3306
```

### 2. 遠程測試
```bash
# 測試資料庫連接
mysql -h YOUR_CLOUD_SQL_IP -u app -p business_exchange
```

## 📝 注意事項

1. **IP 白名單**: 確保您的應用程式 IP 在 Cloud SQL 授權網路中
2. **密碼安全**: 使用強密碼，定期更換
3. **SSL 連接**: 生產環境必須使用 SSL
4. **備份**: 設置自動備份策略
5. **監控**: 監控資料庫性能和連接數

## 🆘 故障排除

### 常見問題

1. **連接被拒絕**
   - 檢查 IP 白名單
   - 確認資料庫用戶權限

2. **認證失敗**
   - 檢查用戶名和密碼
   - 確認用戶主機設置

3. **SSL 錯誤**
   - 確認 SSL 配置
   - 檢查證書有效性

## 📞 支援

如果遇到問題：
1. 檢查 Google Cloud Console 的錯誤日誌
2. 查看 Cloud SQL 實例狀態
3. 參考 Google Cloud 文檔
4. 聯繫 Google Cloud 支援
