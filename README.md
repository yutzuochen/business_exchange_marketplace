# Business Exchange Marketplace

一個類 BizBuySell 的企業互惠平台，使用 Go 技術棧構建。

## 🚀 技術棧

### Backend
- **Go 1.22** - 主要程式語言
- **Gin** - HTTP Web 框架
- **GORM** - ORM 框架
- **MySQL 8** - 主要資料庫
- **Redis** - 快取和會話存儲
- **JWT** - 身份驗證
- **Zap** - 結構化日誌
- **Wire** - 依賴注入

### API
- **REST API** - 標準 RESTful 接口
- **GraphQL** - 使用 gqlgen 的 GraphQL 服務

### Frontend
- **Go Templates** - 服務端渲染
- **Tailwind CSS** - 樣式框架（CDN）

### Infrastructure
- **Docker Compose** - 本地開發環境
- **Makefile** - 構建和部署腳本

## 📁 專案結構

```
business_exchange_marketplace/
├── cmd/                    # 應用程式入口點
│   └── server/            # 主服務器
├── internal/               # 內部包
│   ├── auth/              # 認證相關
│   ├── config/            # 配置管理
│   ├── database/          # 資料庫連接和遷移
│   ├── graphql/           # GraphQL 相關
│   ├── handlers/          # HTTP 處理器
│   ├── logger/            # 日誌配置
│   ├── middleware/        # HTTP 中介層
│   ├── models/            # 資料模型
│   ├── redisclient/       # Redis 客戶端和快取
│   └── router/            # 路由配置
├── templates/              # HTML 模板
├── graph/                  # GraphQL schema 和 resolvers
├── static/                 # 靜態文件
├── uploads/                # 上傳文件
├── docker-compose.yml      # Docker 編排
├── Dockerfile             # 應用程式容器化
├── Makefile               # 構建腳本
├── go.mod                 # Go 模組
├── gqlgen.yml             # GraphQL 生成配置
└── env.example            # 環境變數範例
```

## 🛠️ 本機啟動步驟

### 1. 環境準備

```bash
# 克隆專案
git clone <repository-url>
cd business_exchange_marketplace

# 安裝 Go 1.22+
go version

# 安裝依賴
go mod tidy
```

### 2. 環境配置

```bash
# 複製環境變數範例
cp env.example .env

# 編輯 .env 文件，配置資料庫和 Redis 連接
vim .env
```

### 3. 啟動服務

```bash
# 使用 Docker Compose 啟動所有服務
make docker-up

# 或者分別啟動
docker compose up -d mysql redis
docker compose up -d app
```

### 4. 驗證服務

```bash
# 檢查服務狀態
docker compose ps

# 測試健康檢查
curl http://localhost:8080/healthz

# 訪問應用程式
open http://localhost:8080

# 訪問 Adminer（資料庫管理）
open http://localhost:8081
```

### 5. 開發模式

```bash
# 本地開發（需要本地 MySQL 和 Redis）
go run ./cmd/server

# 或者使用 Makefile
make run
```

## 📋 驗收清單

### 基礎功能
- [x] 專案目錄結構（/cmd, /internal, /pkg）
- [x] Go 模組配置（go.mod）
- [x] 環境變數配置（.env.example）
- [x] Docker Compose 配置
- [x] Makefile 構建腳本

### 資料模型
- [x] 用戶模型（users）
- [x] 刊登模型（listings）
- [x] 圖片模型（images）
- [x] 收藏模型（favorites）
- [x] 訊息模型（messages）
- [x] 交易模型（transactions）
- [x] 資料庫遷移（Auto-migrate）

### API 功能
- [x] REST API 雛形
- [x] GraphQL Schema 雛形
- [x] 用戶註冊/登入
- [x] 刊登 CRUD 操作
- [x] 收藏功能
- [x] 訊息系統

### 中介層
- [x] JWT 認證
- [x] Request ID 追蹤
- [x] 錯誤統一處理
- [x] CORS 配置
- [x] Panic Recovery

### 快取系統
- [x] Redis 快取模組
- [x] 搜尋結果快取
- [x] TTL 配置
- [x] 快取失效策略

### 前端頁面
- [x] 首頁（index.html）
- [x] 註冊頁面（register.html）
- [x] 登入頁面（login.html）
- [x] 儀表板（dashboard.html）
- [x] 市場首頁（market_home.html）
- [x] 刊登詳情（market_listing.html）

### 部署配置
- [x] Dockerfile
- [x] Docker Compose
- [x] 環境變數配置
- [x] 健康檢查端點

## 🔧 常用命令

```bash
# 構建應用程式
make build

# 運行應用程式
make run

# 清理構建文件
make clean

# 更新依賴
make tidy

# 生成 GraphQL 代碼
make gqlgen

# 生成 Wire 依賴注入
make wire

# 啟動 Docker 服務
make docker-up

# 停止 Docker 服務
make docker-down
```

## 🌐 API 端點

### 公開端點
- `GET /` - 首頁
- `GET /market` - 市場首頁
- `GET /market/search` - 搜尋刊登
- `GET /market/listings/:id` - 刊登詳情
- `GET /login` - 登入頁面
- `GET /register` - 註冊頁面
- `GET /healthz` - 健康檢查

### REST API
- `POST /api/v1/auth/register` - 用戶註冊
- `POST /api/v1/auth/login` - 用戶登入
- `GET /api/v1/listings` - 獲取刊登列表
- `GET /api/v1/listings/:id` - 獲取刊登詳情
- `GET /api/v1/categories` - 獲取分類列表

### GraphQL
- `POST /graphql` - GraphQL 查詢
- `GET /playground` - GraphQL Playground

## 🚀 部署到 GCP

### 準備工作
1. 安裝 Google Cloud SDK
2. 配置專案和認證
3. 啟用必要的 API 服務

### 部署步驟
```bash
# 構建容器映像
docker build -t gcr.io/PROJECT_ID/business-exchange .

# 推送到 Google Container Registry
docker push gcr.io/PROJECT_ID/business-exchange

# 部署到 Cloud Run
gcloud run deploy business-exchange \
  --image gcr.io/PROJECT_ID/business-exchange \
  --platform managed \
  --region asia-east1 \
  --allow-unauthenticated
```

## 📝 開發筆記

- 使用 `go mod tidy` 更新依賴
- 使用 `make gqlgen` 重新生成 GraphQL 代碼
- 使用 `make wire` 重新生成依賴注入代碼
- 檢查 `docker-compose.yml` 中的服務健康檢查

## 🤝 貢獻

1. Fork 專案
2. 創建功能分支
3. 提交變更
4. 推送到分支
5. 創建 Pull Request

## 📄 授權

本專案採用 MIT 授權條款。
# Trigger GitHub Actions
# Trigger deployment after fixing GCP_SA_KEY
