# 🚀 CI/CD 設置指南

## 📋 概述

本專案使用 GitHub Actions 實現完整的 CI/CD 流程，包括：
- 🔍 代碼檢查和測試
- 🐳 Docker 映像構建和推送
- 🚀 自動部署到 Google Cloud Run
- 📊 部署後檢查和通知

## 🏗️ 架構

```
GitHub Repository
       ↓
   GitHub Actions
       ↓
   ┌─────────────────┐    ┌─────────────────┐
   │   Staging      │    │   Production    │
   │   Environment  │    │   Environment   │
   └─────────────────┘    └─────────────────┘
       ↓                        ↓
   Cloud Run              Cloud Run
   (Staging)              (Production)
       ↓                        ↓
   Cloud SQL              Cloud SQL
   (Staging DB)           (Production DB)
```

## 🔧 設置步驟

### 1. GitHub Secrets 配置

在 GitHub Repository 中設置以下 Secrets：

#### Google Cloud 認證
```bash
GCP_SA_KEY: 服務帳戶的 JSON 金鑰
```

#### 資料庫配置
```bash
DB_HOST: Cloud SQL 實例 IP
DB_USER: 資料庫用戶名
DB_PASSWORD: 資料庫密碼
DB_NAME: 資料庫名稱
JWT_SECRET: JWT 密鑰
REDIS_ADDR: Redis 地址
```

#### 通知配置 (可選)
```bash
SLACK_WEBHOOK: Slack Webhook URL
```

### 2. 分支策略

```bash
main      → 生產環境部署
develop   → 測試環境部署
feature/* → 功能開發分支
hotfix/*  → 緊急修復分支
```

### 3. 環境保護規則

#### Staging 環境
- 需要 1 個審核
- 自動部署

#### Production 環境
- 需要 2 個審核
- 需要代碼擁有者審核
- 需要通過所有檢查

## 🔄 工作流程

### 1. 代碼檢查和測試 (lint-and-test)
- Go 代碼格式化檢查
- 靜態代碼分析
- 單元測試
- 安全掃描 (Trivy)
- 代碼覆蓋率檢查

### 2. 構建和推送 (build-and-push)
- 多平台 Docker 映像構建
- 推送到 Google Container Registry
- 智能標籤管理
- 構建緩存優化

### 3. 部署到 Staging (deploy-staging)
- 自動部署到測試環境
- 健康檢查
- Slack 通知

### 4. 部署到 Production (deploy-production)
- 手動觸發部署
- 生產環境配置
- 版本標籤管理
- 部署狀態更新

### 5. 部署後檢查 (post-deploy)
- 服務狀態檢查
- 部署狀態更新
- 監控和日誌

## 🚀 觸發條件

### 自動觸發
- **Push to main**: 觸發生產部署
- **Push to develop**: 觸發測試部署
- **Pull Request**: 觸發代碼檢查

### 手動觸發
- **Release**: 發布新版本
- **Workflow Dispatch**: 手動執行

## 📊 監控和通知

### 1. GitHub Actions 監控
- 工作流程執行狀態
- 構建和部署日誌
- 失敗通知

### 2. Slack 通知
- 部署開始通知
- 部署成功/失敗通知
- 環境狀態更新

### 3. Google Cloud 監控
- Cloud Run 服務狀態
- 資料庫連接狀態
- 應用程式日誌

## 🔍 故障排除

### 常見問題

#### 1. 認證失敗
```bash
# 檢查 GCP_SA_KEY 是否正確
# 確認服務帳戶權限
# 檢查專案 ID 是否正確
```

#### 2. 部署失敗
```bash
# 檢查環境變數設置
# 確認 Cloud Run 服務存在
# 檢查網路連接
```

#### 3. 測試失敗
```bash
# 檢查 Go 版本
# 確認依賴完整性
# 檢查測試配置
```

### 調試步驟

1. **檢查 GitHub Actions 日誌**
2. **驗證環境變數**
3. **測試本地構建**
4. **檢查 Google Cloud 權限**

## 📈 最佳實踐

### 1. 代碼質量
- 使用 linting 工具
- 保持測試覆蓋率 > 80%
- 定期安全掃描

### 2. 部署策略
- 藍綠部署
- 滾動更新
- 自動回滾

### 3. 監控和警報
- 設置關鍵指標監控
- 配置自動警報
- 定期健康檢查

## 🔐 安全考慮

### 1. 密鑰管理
- 使用 GitHub Secrets
- 定期輪換密鑰
- 最小權限原則

### 2. 網路安全
- 使用私有網路
- 啟用 SSL/TLS
- 配置防火牆規則

### 3. 代碼安全
- 依賴漏洞掃描
- 容器安全掃描
- 代碼簽名驗證

## 📚 參考資源

- [GitHub Actions 文檔](https://docs.github.com/en/actions)
- [Google Cloud Run 文檔](https://cloud.google.com/run/docs)
- [Go 測試指南](https://golang.org/doc/tutorial/testing)
- [Docker 最佳實踐](https://docs.docker.com/develop/dev-best-practices/)

## 🆘 支援

如果遇到問題：
1. 檢查 GitHub Actions 日誌
2. 查看 Google Cloud Console
3. 參考故障排除指南
4. 聯繫開發團隊
