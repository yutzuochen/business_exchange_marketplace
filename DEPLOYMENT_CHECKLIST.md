# Cloud Run 部署檢查清單

## ✅ 已完成的配置

### 1. 服務器綁定配置
- ✅ **服務器地址**: `":" + cfg.AppPort` (綁定到 `0.0.0.0:$PORT`)
- ✅ **PORT 環境變數**: Cloud Run 會自動設置，代碼已正確讀取
- ✅ **位置**: `cmd/server/main.go:93`

### 2. 數據庫連接配置
- ✅ **Unix Socket 支持**: 已在 `config.go` 中實現
- ✅ **自動檢測**: 當 `DB_HOST` 以 `/` 開頭時使用 Unix Socket
- ✅ **連接字符串**: 正確格式化 Cloud SQL Unix Socket 連接
- ✅ **位置**: `internal/config/config.go:154-158`

### 3. 部署腳本配置
- ✅ **deploy-to-cloud.sh**: 已更新使用 Unix Socket
- ✅ **deploy.sh**: 已更新後端部署配置
- ✅ **環境變數**: 正確設置 Cloud SQL 連接參數

### 4. 環境變數設置
- ✅ **DB_HOST**: `/cloudsql/businessexchange-468413:us-central1:trade-sql`
- ✅ **Cloud SQL 實例**: `--add-cloudsql-instances` 已配置
- ✅ **生產配置**: 創建了 `env.production` 文件

## 🚀 部署步驟

### 1. 準備部署
```bash
cd /home/mason/Documents/bex567/business_exchange_marketplace

# 確保已登入 Google Cloud
gcloud auth login
gcloud auth configure-docker

# 設置專案
gcloud config set project businessexchange-468413
```

### 2. 檢查 Cloud SQL 實例
```bash
# 確認實例正在運行
gcloud sql instances describe trade-sql --project=businessexchange-468413

# 如果需要啟動實例
gcloud sql instances patch trade-sql --activation-policy=ALWAYS
```

### 3. 運行數據庫遷移（首次部署）
```bash
./run-migrations-cloud.sh
```

### 4. 部署應用
```bash
# 選項 A: 簡單部署（推薦）
./deploy-to-cloud.sh

# 選項 B: 完整部署（包含前端）
./deploy.sh
```

### 5. 驗證部署
```bash
# 腳本會自動測試健康檢查
# 手動測試：
curl https://your-service-url/health

# 查看日誌
gcloud logs read --service=business-exchange --limit=20
```

## 🔧 關鍵配置說明

### 服務器綁定
- Cloud Run 會設置 `PORT` 環境變數（通常是 8080）
- 服務器綁定到 `0.0.0.0:$PORT`，允許接收來自任何 IP 的請求
- 這是 Cloud Run 的標準要求

### 數據庫連接
- 使用 Unix Socket: `/cloudsql/PROJECT_ID:REGION:INSTANCE_NAME`
- Cloud Run 會自動注入 Cloud SQL Proxy
- 不需要 IP 地址或端口號
- 連接更安全且性能更好

### 環境變數
```bash
APP_ENV=production
DB_HOST=/cloudsql/businessexchange-468413:us-central1:trade-sql
DB_USER=app
DB_PASSWORD=app_password
DB_NAME=business_exchange
```

## ⚠️  注意事項

1. **JWT Secret**: 請在生產環境中更換為強密碼
2. **CORS 設置**: 生產環境應設置具體的允許域名
3. **數據庫密碼**: 考慮使用 Google Secret Manager
4. **健康檢查**: 確保你的應用有 `/health` 端點

## 🔍 故障排除

### 如果部署失敗
```bash
# 查看部署日誌
gcloud logs read --service=business-exchange --limit=50

# 檢查服務狀態
gcloud run services describe business-exchange --region=us-central1
```

### 如果數據庫連接失敗
1. 確認 Cloud SQL 實例正在運行
2. 檢查實例名稱是否正確
3. 確認數據庫用戶和密碼
4. 查看應用日誌中的錯誤信息

### 如果服務無法訪問
1. 確認服務已設置為 `--allow-unauthenticated`
2. 檢查服務是否綁定到正確的端口
3. 驗證健康檢查端點

## 📝 部署後任務

1. 測試所有 API 端點
2. 設置監控和告警
3. 配置自定義域名（可選）
4. 設置 SSL 證書（自動）
5. 配置備份策略
